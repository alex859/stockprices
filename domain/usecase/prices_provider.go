package usecase

import (
	"org.alex859/stockprices/domain/entity"
)

type (
	// CurrentPriceProvider returns the current price for the given Ticker.
	// If nothing can be found, return an ErrNothingFound error.
	CurrentPriceProvider interface {
		GetCurrentPrice(ticker entity.Ticker) (entity.CurrentPrice, error)
	}

	// HistoricalPricesProvider returns price history for a given ticker in the given date interval.
	// Returned prices are chronologically ordered.
	// If nothing can be found, return an ErrNothingFound error.
	HistoricalPricesProvider interface {
		GetHistoricalPrices(ticker entity.Ticker, dateInterval entity.DateInterval) (entity.PriceHistory, error)
	}
)

