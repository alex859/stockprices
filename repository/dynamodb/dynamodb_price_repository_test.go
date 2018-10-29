package dynamodb

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"context"
	dynamodb2 "github.com/PolarGeospatialCenter/dockertest/pkg/dynamodb"
	"github.com/stretchr/testify/assert"
	"org.alex859/stockprices/domain"
	"time"
	"github.com/pkg/errors"
)

func Test_dynamoDbPricesRepository_createTableWriteAndRead(t *testing.T) {
	ctx := context.Background()
	instance, err := dynamodb2.Run(ctx)
	if err != nil {
		t.Fatalf("Unable to start dynamodb: %v", err)
	}
	defer instance.Stop(ctx)

	repository := NewDynamoDbPriceRepository(instance.Config())

	if assert.NoError(t, repository.createTable()) {
		s := session.Must(session.NewSession(instance.Config()))

		cli := dynamodb.New(s)
		out, err := cli.ListTables(&dynamodb.ListTablesInput{})
		if err != nil {
			t.Fatalf("Unable to list tables: %v", err)
		}
		assert.Equal(t, 1, len(out.TableNames))
		assert.Equal(t, tableName, *out.TableNames[0])

		ticker := domain.Ticker{Market:"LON", Symbol:"ANP"}

		var may5, _ = time.Parse("02-01-2006", "05-05-2018")
		var may9, _ = time.Parse("02-01-2006", "09-05-2018")
		ppMay5 := domain.PricePoint{Price: 123, Time: may5}
		ppMay9 := domain.PricePoint{Price: 123, Time: may9}

		period := domain.FiveDays

		// we still haven't written anything
		_, err = repository.Read(ticker, period)
		assert.Equal(t, domain.NewErrNothingFound(ticker), errors.Cause(err))

		// first write in 2018
		var firstWriteTime, _ = time.Parse("02-01-2006", "09-05-2018")
		now = func() time.Time {
			return firstWriteTime
		}
		prices := domain.PriceList{ppMay5, ppMay9}
		priceHistory := domain.PriceHistory{TickerInfo: domain.TickerInfo{Ticker:ticker}, Prices:prices}
		repository.Save(ticker, period, priceHistory)

		result, err := repository.Read(ticker, period)
		if assert.NoError(t, err) {
			assert.Equal(t, priceHistory, result)
		}

		// after 1 year..
		var afterOneYear, _ = time.Parse("02-01-2006", "09-05-2019")
		now = func() time.Time {
			return afterOneYear
		}
		_, err = repository.Read(ticker, period)
		assert.Equal(t, domain.NewErrNothingFound(ticker), err)
	}
}

