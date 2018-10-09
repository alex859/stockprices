package main

import (
	"net/http"
	"org.alex859/stockprices/price/interfaces"
	"org.alex859/stockprices/domain"
	"github.com/aws/aws-lambda-go/events"
	"time"
	"encoding/json"
	"log"
	"strings"
	"org.alex859/stockprices/price/usecases"
	"fmt"
	"os"
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
	//lambda.Start(Handler)
	configuration := ReadConfig()
	priceProvider := interfaces.NewGooglePricesProvider(interfaces.NewGooglePrices(http.DefaultClient))
	from, _ := time.Parse("02-01-2006", "04-09-2008")
	to, _ := time.Parse("02-01-2006", "13-09-2018")
	uc := usecases.NewGetPricesUseCase(priceProvider, configuration.NumPriceProviderWorkers)

	start := time.Now()

	result, err := uc.GetMultiplePrices([]domain.Ticker{
		{Market: "LON", Symbol: "ANP"},
		{Market: "LON", Symbol: "TEF"},
		{Market: "LON", Symbol: "BKG"},
		{Market: "LON", Symbol: "NG"},
		{Market: "LON", Symbol: "SMDS"},
		{Market: "LON", Symbol: "AHT"},
		{Market: "LON", Symbol: "RFX"},
		{Market: "LON", Symbol: "BMY"},
		{Market: "LON", Symbol: "HSBC"},
		{Market: "LON", Symbol: "GATC"},
		{Market: "LON", Symbol: "ANP"},
		{Market: "LON", Symbol: "TEF"},
		{Market: "LON", Symbol: "BKG"},
		{Market: "LON", Symbol: "NG"},
		{Market: "LON", Symbol: "SMDS"},
		{Market: "LON", Symbol: "AHT"},
		{Market: "LON", Symbol: "RFX"},
		{Market: "LON", Symbol: "BMY"},
		{Market: "LON", Symbol: "HSBC"},
		{Market: "LON", Symbol: "GATC"},
	}, from, to)

	elapsed := time.Since(start)
	log.Printf("It took %s", elapsed)

	if err != nil {
		log.Println(err.Error())
	}

	resultMap := make(map[string]*domain.PriceHistory, len(result))
	for key, value := range result {
		resultMap[fmt.Sprintf("%s:%s", key.Market, key.Symbol)] = value
	}

	s, err := json.Marshal(resultMap)

	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(s))
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