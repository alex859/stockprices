package usecase

import (
	"org.alex859/stockprices/domain/entity"
	"log"
	"github.com/pkg/errors"
)

type getHistoricalPricesUseCase struct {
	priceProvider HistoricalPricesProvider
	numWorkers    int
}
func NewGetHistoricalPricesUseCase(dataProvider HistoricalPricesProvider, numWorkers int) *getHistoricalPricesUseCase {
	return &getHistoricalPricesUseCase{priceProvider: dataProvider, numWorkers:numWorkers}
}

type historicalPricesResultErrorChannel struct {
	result entity.PriceHistory
	err error
}

func (useCase *getHistoricalPricesUseCase) GetHistoricalPrices(tickers []entity.Ticker, interval entity.DateInterval) (map[string]entity.PriceHistory, error) {
	n := len(tickers)
	if n == 0 {
		return map[string]entity.PriceHistory{}, nil
	}

	resultsChannel := make(chan historicalPricesResultErrorChannel, n)
	tickersChannel := make(chan entity.Ticker, n)

	for w := 1; w <= useCase.numWorkers; w++ {
		go historicalPricesProviderWorker(useCase.priceProvider, tickersChannel, resultsChannel, interval)
	}

	for _, ticker := range tickers {
		tickersChannel <- ticker
	}
	close(tickersChannel)

	result := map[string]entity.PriceHistory{}
	for i := 0; i < n; i++ {
		r := <-resultsChannel
		if r.err == nil {
			result[r.result.Ticker.String()] = r.result
		}
	}

	if len(result) == 0 {
		return nil, errors.New("unable to fetch stock prices")
	}

	return result, nil
}

func historicalPricesProviderWorker(priceProvider HistoricalPricesProvider, tickersChannel <-chan entity.Ticker, ch chan<- historicalPricesResultErrorChannel, interval entity.DateInterval) {
	for ticker := range tickersChannel {
		history, err := priceProvider.GetHistoricalPrices(ticker, interval)
		if err != nil {
			log.Printf("An error occured while fetching prices for ticker: %s, interval: %s. Error: %+v", ticker.String(), interval, err)
		}
		ch <- historicalPricesResultErrorChannel{result:history, err:err}
	}
}