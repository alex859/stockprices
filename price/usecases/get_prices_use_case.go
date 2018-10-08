package usecases

import (
	"org.alex859/stockprices/price"
	"org.alex859/stockprices/domain"
	"log"
	"errors"
	"time"
)

type getPricesUseCase struct {
	priceProvider PriceProvider
}
func NewGetPricesUseCase(dataProvider PriceProvider) price.GetPricesUseCase {
	return &getPricesUseCase{priceProvider: dataProvider}
}

type resultError struct {
	result *domain.PriceHistory
	err error
}

func (useCase *getPricesUseCase) GetMultiplePrices(tickers []domain.Ticker, from time.Time, to time.Time) (map[domain.Ticker]*domain.PriceHistory, error) {
	n := len(tickers)
	if n == 0 {
		return map[domain.Ticker]*domain.PriceHistory{}, nil
	}

	resultsChannel := make(chan resultError, n)

	for _, ticker := range tickers {
		go func(ticker domain.Ticker, from time.Time, to time.Time, ch chan<- resultError) {
			history, err := useCase.priceProvider.FetchPrices(ticker, from, to)
			if err != nil {
				log.Printf("An error occured while fetching prices for ticker: %s, from: %s, to: %s. Error: %s", ticker.String(), from, to, err)
			}
			ch <- resultError{result:history, err:err}
		}(ticker, from, to, resultsChannel)
	}

	result := map[domain.Ticker]*domain.PriceHistory{}
	for i := 0; i < n; i++ {
		r := <-resultsChannel
		if r.err == nil {
			result[r.result.Ticker] = r.result
		}
	}

	if len(result) == 0 {
		return nil, errors.New("unable to fetch stock prices")
	}

	return result, nil
}