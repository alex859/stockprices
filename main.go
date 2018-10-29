package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"org.alex859/stockprices/domain"
	"org.alex859/stockprices/pricesprovider"
	"org.alex859/stockprices/usecases"
	"org.alex859/stockprices/googlefinance"
	"fmt"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//configuration := ReadConfig()
	googleFinancePriceFetcher := googlefinance.NewDefaultGooglePricesFetcher(http.DefaultClient)
	pricesProvider := pricesprovider.NewGoogleFinancePricesProvider(googleFinancePriceFetcher, googlefinance.NewGoogleFinanceResponseConverter())
	useCase := usecases.NewGetHistoricalPricesUseCase(pricesProvider, 5)
	symbols := strings.Split(request.QueryStringParameters["symbols"], ",")
	fromStr := request.QueryStringParameters["from"]
	toStr := request.QueryStringParameters["to"]
	from, err := time.Parse("02-01-2006", fromStr)
	to, err := time.Parse("02-01-2006", toStr)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	var tickers []domain.Ticker
	for _, symbol := range symbols {
		split := strings.Split(symbol, ":")
		if len(split) != 2 {
			return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Unrecognized symbol %s", symbol), StatusCode: 400}, nil
		}
		tickers = append(tickers, domain.Ticker{Market:split[0], Symbol:split[1]})
	}


	interval, err := domain.NewDateInterval(from, to)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}
	result, err := useCase.GetHistoricalPrices(tickers, interval)

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

	//googleFinancePriceFetcher := googlefinance.NewDefaultGooglePricesFetcher(http.DefaultClient)
	//pricesProvider := pricesprovider.NewGoogleFinancePricesProvider(googleFinancePriceFetcher, googlefinance.NewGoogleFinanceResponseConverter())
	//useCase := usecases.NewGetHistoricalPricesUseCase(pricesProvider, configuration.NumPriceProviderWorkers)
	//from, _ := time.Parse("02-01-2006", "04-09-2000")
	//to, _ := time.Parse("02-01-2006", "13-09-2018")
	//
	//start := time.Now()
	//
	//interval, err := domain.NewDateInterval(from, to)
	//if err != nil {
	//	return
	//}
	//result, err := useCase.GetHistoricalPrices([]domain.Ticker{
	//	{Market: "LON", Symbol: "HSBSSC"},}, interval)
	//
	//elapsed := time.Since(start)
	//log.Printf("It took %s", elapsed)
	//
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//
	//resultMap := make(map[string]domain.PriceHistory, len(result))
	//for key, value := range result {
	//	resultMap[fmt.Sprintf("%s:%s", key.Market, key.Symbol)] = value
	//}
	//
	//s, err := json.Marshal(resultMap)
	//
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//log.Println(string(s))
}

type Configuration struct {
	NumPriceProviderWorkers int
}

func ReadConfig() Configuration {
	file, err := os.Open("config/config.dev.json")
	if err != nil {
		log.Fatal(err)
	}
	configuration := Configuration{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	return configuration
}