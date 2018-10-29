package usecases

import (
	"org.alex859/stockprices/domain"
)

type (
	// Return the current price for the given Ticker.
	// If nothing can be found, return an ErrNothingFound error.
	CurrentPriceProvider interface {
		GetCurrentPrice(ticker domain.Ticker) (domain.CurrentPrice, error)
	}

	// Returns price history for a given ticker in the given date interval.
	// Returned prices are chronologically ordered.
	// If nothing can be found, return an ErrNothingFound error.
	HistoricalPricesProvider interface {
		GetHistoricalPrices(ticker domain.Ticker, dateInterval domain.DateInterval) (domain.PriceHistory, error)
	}

	PricesProvider interface {
		CurrentPriceProvider
		HistoricalPricesProvider
	}
)

