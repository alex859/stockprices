package usecases

import (
	"org.alex859/stockprices/price"
	"org.alex859/stockprices/domain"
	"log"
	"errors"
	"time"
)

type getPricesUseCase struct {
	priceProvider PricesProvider
	numWorkers    int
}
func NewGetPricesUseCase(dataProvider PricesProvider, numWorkers int) price.GetPricesUseCase {
	return &getPricesUseCase{priceProvider: dataProvider, numWorkers:numWorkers}
}

type resultErrorChannel struct {
	result *domain.PriceHistory
	err error
}

func (useCase *getPricesUseCase) GetMultiplePrices(tickers []domain.Ticker, from time.Time, to time.Time) (map[domain.Ticker]*domain.PriceHistory, error) {
	n := len(tickers)
	if n == 0 {
		return map[domain.Ticker]*domain.PriceHistory{}, nil
	}

	resultsChannel := make(chan resultErrorChannel, n)
	tickersChannel := make(chan domain.Ticker, n)

	for w := 1; w <= useCase.numWorkers; w++ {
		go priceProviderWorker(useCase.priceProvider, tickersChannel, resultsChannel, from, to)
	}

	for _, ticker := range tickers {
		tickersChannel <- ticker
	}
	close(tickersChannel)

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

func priceProviderWorker(priceProvider PricesProvider, tickersChannel <-chan domain.Ticker, ch chan<- resultErrorChannel, from time.Time, to time.Time, ) {
	for ticker := range tickersChannel {
		history, err := priceProvider.FetchPrices(ticker, from, to)
		if err != nil {
			log.Printf("An error occured while fetching prices for ticker: %s, from: %s, to: %s. Error: %s", ticker.String(), from, to, err)
		}
		ch <- resultErrorChannel{result:history, err:err}
	}
}