package main

import (
	"github.com/aws/aws-lambda-go/events"
	"org.alex859/stockprices/googlefinance"
	"net/http"
	"org.alex859/stockprices/pricesprovider"
	"org.alex859/stockprices/usecases"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"org.alex859/stockprices/handlers"
)

func HistoricalPricesHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//configuration := ReadConfig()
	googleFinancePriceFetcher := googlefinance.NewDefaultGooglePricesFetcher(http.DefaultClient)
	pricesProvider := pricesprovider.NewGoogleFinancePricesProvider(googleFinancePriceFetcher, googlefinance.NewGoogleFinanceResponseConverter())
	useCase := usecases.NewGetHistoricalPricesUseCase(pricesProvider, 5)

	tickerSlice, err := handlers.Tickers(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	interval, err := handlers.Interval(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	result, err := useCase.GetHistoricalPrices(tickerSlice, interval)

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
	lambda.Start(HistoricalPricesHandler)
}
