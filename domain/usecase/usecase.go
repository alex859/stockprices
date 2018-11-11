package usecase

import (
	"org.alex859/stockprices/domain/entity"
)

type (
	// GetHistoricalPricesUseCase prices history of the given stocks in the give time interval.
	GetHistoricalPricesUseCase interface {
		GetHistoricalPrices(tickers []entity.Ticker, interval entity.DateInterval) (map[string]entity.PriceHistory, error)
	}

	// GetCurrentPricesUseCase prices history of the given stocks in the give time interval.
	GetCurrentPricesUseCase interface {
		GetCurrentPrices(tickers []entity.Ticker) (map[string]entity.CurrentPrice, error)
	}
)
