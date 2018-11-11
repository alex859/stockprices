package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"org.alex859/stockprices/data/googlefinance"
	"org.alex859/stockprices/domain/usecase"
	"org.alex859/stockprices/presentation/handlers"
)

func CurrentPricesHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//configuration := ReadConfig()
	googleFinancePriceFetcher := googlefinance.NewDefaultPricesFetcher(http.DefaultClient)
	pricesProvider := googlefinance.NewGoogleFinancePricesProvider(googleFinancePriceFetcher, googlefinance.NewGoogleFinanceResponseConverter())
	useCase := usecase.NewGetCurrentPricesUseCase(pricesProvider, 5)

	tickerSlice, err := handlers.Tickers(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	result, err := useCase.GetCurrentPrices(tickerSlice)
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
	lambda.Start(CurrentPricesHandler)
}
