package usecases

import (
	"org.alex859/stockprices/domain"
)

type (
	// Get prices history of the given stocks in the give time interval.
	GetHistoricalPricesUseCase interface {
		GetHistoricalPrices(tickers []domain.Ticker, interval domain.DateInterval) (map[string]domain.PriceHistory, error)
	}

	// Get prices history of the given stocks in the give time interval.
	GetCurrentPricesUseCase interface {
		GetCurrentPrices(tickers []domain.Ticker) (map[string]domain.CurrentPrice, error)
	}
)
