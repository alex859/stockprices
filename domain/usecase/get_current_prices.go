package usecase

import (
	"org.alex859/stockprices/domain/entity"
	"log"
	"errors"
)

type getCurrentPricesUseCase struct {
	priceProvider CurrentPriceProvider
	numWorkers    int
}
func NewGetCurrentPricesUseCase(dataProvider CurrentPriceProvider, numWorkers int) *getCurrentPricesUseCase {
	return &getCurrentPricesUseCase{priceProvider: dataProvider, numWorkers:numWorkers}
}

type currentPricesResultErrorChannel struct {
	result entity.CurrentPrice
	err error
}

func (useCase *getCurrentPricesUseCase) GetCurrentPrices(tickers []entity.Ticker) (map[string]entity.CurrentPrice, error) {
	n := len(tickers)
	if n == 0 {
		return map[string]entity.CurrentPrice{}, nil
	}

	resultsChannel := make(chan currentPricesResultErrorChannel, n)
	tickersChannel := make(chan entity.Ticker, n)

	for w := 1; w <= useCase.numWorkers; w++ {
		go currentPriceProviderWorker(useCase.priceProvider, tickersChannel, resultsChannel)
	}

	for _, ticker := range tickers {
		tickersChannel <- ticker
	}
	close(tickersChannel)

	result := map[string]entity.CurrentPrice{}
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

func currentPriceProviderWorker(priceProvider CurrentPriceProvider, tickersChannel <-chan entity.Ticker, ch chan<- currentPricesResultErrorChannel) {
	for ticker := range tickersChannel {
		history, err := priceProvider.GetCurrentPrice(ticker)
		if err != nil {
			log.Printf("An error occured while fetching prices for ticker: %s. Error: %+v", ticker.String(), err)
		}
		ch <- currentPricesResultErrorChannel{result:history, err:err}
	}
}