package interfaces

import (
	"org.alex859/stockprices/price/usecases"
	"org.alex859/stockprices/domain"
	"time"
	"errors"
	"log"
)

type googlePricesProvider struct {
	googlePricesClient *googlePricesClient
}

func NewGooglePricesProvider(client *googlePricesClient) usecases.PricesProvider {
	return &googlePricesProvider{client}
}

func (gpp *googlePricesProvider) FetchPrices(ticker domain.Ticker, from time.Time, to time.Time) (*domain.PriceHistory, error) {

	period, err := timeToPeriod(from, to)
	if err != nil {
		log.Printf("Error converting time interval to period: %s", err)
		return nil, err
	}

	priceHistory, err:= gpp.googlePricesClient.FetchPrices(ticker, period)
	if err != nil {
		log.Printf("getting prices from Google: %s", err)
		return nil, err
	}

	priceHistory.Prices = filterPriceList(priceHistory.Prices, from, to)

	return priceHistory, nil
}


func timeToPeriod(from time.Time, to time.Time) (Period, error) {

	if to.Before(from) {
		return Unknown, errors.New("parameter to must be after parameter from")
	}

	today := time.Now()
	sameDay := from.Year() == today.Year() && from.Month() == today.Month() && from.Day() == today.Day()
	diff := time.Since(from)

	day := time.Hour * 24
	month := day * 30
	year := day * 365
	switch {
	case sameDay:
		return OneDay, nil
	case diff >= 5*year:
		return Max, nil
	case diff >= 1*year:
		return FiveYears, nil
	case diff >= 6*month:
		return OneYear, nil
	case diff >= 1*month:
		return SixMonth, nil
	case diff >= 5*day:
		return OneMonth, nil
	default:
		return FiveDays, nil
	}
}


func filterPriceList(list domain.PriceList, from time.Time, to time.Time) domain.PriceList {
	if to.Before(from) {
		return domain.PriceList{}
	}
	var result = domain.PriceList{}
	for _, price := range list {
		if (price.Time.After(from) || price.Time.Equal(from)) && (price.Time.Before(to) || price.Time.Equal(to)) {
			result = append(result, price)
		}
	}
	return result
}
