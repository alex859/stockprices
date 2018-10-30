package handlers

import (
	"time"
	"github.com/aws/aws-lambda-go/events"
	"org.alex859/stockprices/domain"
	"github.com/pkg/errors"
	"strings"
)

var now = time.Now

var dateLayout = "02-01-2006"
var fromDateParam = "from"
var toDateParam = "to"
var tickersParam = "tickers"

// From date coming form the request is required. Format DD-MM-YYYY
func FromDate(request events.APIGatewayProxyRequest) (time.Time, error) {
	if str, ok := request.QueryStringParameters[fromDateParam]; ok {
		result, err := time.Parse(dateLayout, str)
		return result, errors.Wrap(err, "Invalid from date parameter")
	}

	// TODO custom errors to be mapped to error codes
	return time.Now(), errors.New("Missing from date parameter")
}

// To date is optional. Defaults to current time. Format DD-MM-YYYY
func ToDate(request events.APIGatewayProxyRequest) (time.Time, error) {
	if str, ok := request.QueryStringParameters[toDateParam]; ok {
		result, err := time.Parse(dateLayout, str)
		return result, errors.Wrap(err, "Invalid to date parameter")
	}

	// TODO custom errors to be mapped to error codes
	return now(), nil
}

func Interval(request events.APIGatewayProxyRequest) (result domain.DateInterval, err error) {
	from, err := FromDate(request)
	if err != nil {
		return
	}
	to, err := ToDate(request)
	if err != nil {
		return
	}
	result, err = domain.NewDateInterval(from, to)
	return result, errors.Wrap(err, "Invalid interval")
}

// At least one ticker has to be present in the form of "MARKET1:SYMBOL1[,MARKET2:SYMBOL2]"
func Tickers(request events.APIGatewayProxyRequest) (result []domain.Ticker, err error) {
	if tickerStr, ok := request.QueryStringParameters[tickersParam]; ok {
		// ignore whitespaces
		tickerStr = strings.Replace(tickerStr, " ", "", -1)
		for _, token := range strings.Split(tickerStr, ",") {
			tickerComponents := strings.Split(token, ":")
			if len(tickerComponents) == 2 {
				result = append(result, domain.Ticker{Market:tickerComponents[0], Symbol:tickerComponents[1]})
			}
		}

		if len(result) == 0 {
			err = errors.New("No valid tickers found")
		}
		return
	}

	err = errors.New("Missing parameter tickers")
	return
}


