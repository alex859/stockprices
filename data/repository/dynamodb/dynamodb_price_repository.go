package dynamodb

import (
	"org.alex859/stockprices/data/googlefinance"
	"org.alex859/stockprices/domain/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"fmt"
	"log"
	"strconv"
	"github.com/pkg/errors"
)

type dynamoDbPricesRepository struct {
	ddb *dynamodb.DynamoDB
}

func NewDynamoDbPriceRepository(config *aws.Config) *dynamoDbPricesRepository {
	s := session.Must(session.NewSession(config))
	return &dynamoDbPricesRepository{ddb: dynamodb.New(s)}
}

var now = time.Now
const tableName = "PriceHistory"
const keyName = "tickerPeriod"

func (r *dynamoDbPricesRepository) createTable() error {

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(keyName),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(keyName),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := r.ddb.CreateTable(input)
	if err != nil {
		return errors.Wrap(err, "unable to create table")
	}

	return nil
}


func (r *dynamoDbPricesRepository) Read(ticker entity.Ticker, period googlefinance.Period) (entity.PriceHistory, error) {
	var priceHistory entity.PriceHistory
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key(ticker, period)),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := r.ddb.GetItem(input)
	if err != nil {
		log.Println(err.Error())
		return priceHistory, errors.Wrap(err, "unable to get item from DynamoDB")
	}

	if result.Item == nil {
		return priceHistory, errors.Wrap(entity.NewErrNothingFound(ticker), "nothing found")
	}

	priceHistoryDDB := PriceHistoryDDB{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &priceHistoryDDB)
	if err != nil {
		log.Println(err.Error())
		return priceHistory, errors.Wrap(err, "unable to read item")
	}

	if priceHistoryDDB.ExpirationDate < now().Unix() {
		log.Println("Found expired item")
		return priceHistory, entity.NewErrNothingFound(ticker)
	}

	return priceHistoryDDB.PriceHistory, nil
}

func (r *dynamoDbPricesRepository) Save(ticker entity.Ticker, period googlefinance.Period, prices entity.PriceHistory) error {
	converted, err := toDDB(prices, period)
	item, err := dynamodbattribute.MarshalMap(converted)
	if err != nil {
		message := fmt.Sprintf("failed to marshal price history, %v", err)
		return errors.Wrap(err, message)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = r.ddb.PutItem(input)
	if err != nil {
		return errors.Wrap(err, "error writing to DynamoDB")
	}

	return nil
}

func key(ticker entity.Ticker, period googlefinance.Period) string {
	return fmt.Sprintf("%v_%v_%v", ticker.Market, ticker.Symbol, period.Value().Str)
}

func toDDB(priceHistory entity.PriceHistory, period googlefinance.Period) (PriceHistoryDDB, error) {
	interval, err := strconv.ParseInt(period.Value().Interval, 10, 64)
	if err != nil {
		return PriceHistoryDDB{}, errors.Wrap(err, "erro converting to DDB item")
	}
	return PriceHistoryDDB{
		priceHistory,
		key(priceHistory.Ticker, period),
		 now().Add(time.Duration(interval) * time.Second).Unix(),
	}, nil
}

type PriceHistoryDDB struct {
	entity.PriceHistory
	Key string `json:"tickerPeriod"`
	ExpirationDate int64
}