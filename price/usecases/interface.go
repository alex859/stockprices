package usecases

import (
	"org.alex859/stockprices/domain"
	"time"
)

type PriceProvider interface {
	FetchPrices(ticker domain.Ticker, from time.Time, to time.Time) (*domain.PriceHistory, error)
}

type PriceSaver interface {
	SavePrices(ticker domain.Ticker, prices *domain.PriceHistory) error
}