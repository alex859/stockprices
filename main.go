package main

import (
	"net/http"
	"org.alex859/stockprices/price/interfaces"
	"org.alex859/stockprices/domain"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"time"
	"encoding/json"
	"log"
	"strings"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	priceProvider := interfaces.NewGooglePricesProvider(interfaces.NewGooglePrices(http.DefaultClient))
	from, _ := time.Parse("02-01-2006", "04-09-2018")
	to, _ := time.Parse("02-01-2006", "13-09-2018")
	symbol := request.QueryStringParameters["symbol"]
	split := strings.Split(symbol, ":")
	if len(split) != 2 {
		return events.APIGatewayProxyResponse{Body: "Unrecognized symbol", StatusCode: 400}, nil
	}

	result, err := priceProvider.FetchPrices(domain.Ticker{Market:split[0], Symbol:split[1]}, from, to)

	log.Printf("You want to know prices!")
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	s, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	return events.APIGatewayProxyResponse{Body: string(s), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
