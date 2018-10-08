package price

import (
	"org.alex859/stockprices/domain"
	"time"
)


type GetPricesUseCase interface {
	GetMultiplePrices(tickers []domain.Ticker, from time.Time, to time.Time) (map[domain.Ticker]*domain.PriceHistory, error)
}

