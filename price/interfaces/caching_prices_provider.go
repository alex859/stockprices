package interfaces

import (
	"org.alex859/stockprices/price/usecases"
	"org.alex859/stockprices/domain"
	"time"
)

type cachingPricesProvider struct {
	pricesProvider usecases.PricesProvider
	cachedPricesProvider usecases.PricesProvider
	pricesSaver usecases.PricesSaver
}

func NewCachingPricesProvider(pricesProvider usecases.PricesProvider, cachedPricesProvider usecases.PricesProvider, pricesSaver usecases.PricesSaver) usecases.PricesProvider {
	return &cachingPricesProvider{pricesProvider:pricesProvider, cachedPricesProvider:cachedPricesProvider, pricesSaver: pricesSaver}
}

func (cpp *cachingPricesProvider) FetchPrices(ticker domain.Ticker, from time.Time, to time.Time) (*domain.PriceHistory, error) {
	return nil, nil
}
