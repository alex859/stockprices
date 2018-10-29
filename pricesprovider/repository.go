package pricesprovider

import (
	"org.alex859/stockprices/domain"
)

type (
	PricesReader interface {
		Read(ticker domain.Ticker, period domain.Period) (domain.PriceHistory, error)
	}

	PricesSaver interface {
		Save(ticker domain.Ticker, period domain.Period, prices domain.PriceHistory) error
	}

	PricesRepository interface {
		PricesReader
		PricesSaver
	}
)
