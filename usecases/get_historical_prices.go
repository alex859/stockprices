package usecases

import (
	"org.alex859/stockprices/domain"
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
	result domain.PriceHistory
	err error
}

func (useCase *getHistoricalPricesUseCase) GetHistoricalPrices(tickers []domain.Ticker, interval domain.DateInterval) (map[string]domain.PriceHistory, error) {
	n := len(tickers)
	if n == 0 {
		return map[string]domain.PriceHistory{}, nil
	}

	resultsChannel := make(chan historicalPricesResultErrorChannel, n)
	tickersChannel := make(chan domain.Ticker, n)

	for w := 1; w <= useCase.numWorkers; w++ {
		go historicalPricesProviderWorker(useCase.priceProvider, tickersChannel, resultsChannel, interval)
	}

	for _, ticker := range tickers {
		tickersChannel <- ticker
	}
	close(tickersChannel)

	result := map[string]domain.PriceHistory{}
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

func historicalPricesProviderWorker(priceProvider HistoricalPricesProvider, tickersChannel <-chan domain.Ticker, ch chan<- historicalPricesResultErrorChannel, interval domain.DateInterval) {
	for ticker := range tickersChannel {
		history, err := priceProvider.GetHistoricalPrices(ticker, interval)
		if err != nil {
			log.Printf("An error occured while fetching prices for ticker: %s, interval: %s. Error: %+v", ticker.String(), interval, err)
		}
		ch <- historicalPricesResultErrorChannel{result:history, err:err}
	}
}