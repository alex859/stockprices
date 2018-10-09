package usecases

import (
	"org.alex859/stockprices/domain"
	"time"
)

type PricesProvider interface {
	FetchPrices(ticker domain.Ticker, from time.Time, to time.Time) (*domain.PriceHistory, error)
}

type PricesSaver interface {
	SavePrices(ticker domain.Ticker, prices *domain.PriceHistory) error
}