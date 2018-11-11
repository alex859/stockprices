package pricesprovider

import (
	"org.alex859/stockprices/data/googlefinance"
	"org.alex859/stockprices/domain/entity"
)

type (
	PricesReader interface {
		Read(ticker entity.Ticker, period googlefinance.Period) (entity.PriceHistory, error)
	}

	PricesSaver interface {
		Save(ticker entity.Ticker, period googlefinance.Period, prices entity.PriceHistory) error
	}

	PricesRepository interface {
		PricesReader
		PricesSaver
	}
)
