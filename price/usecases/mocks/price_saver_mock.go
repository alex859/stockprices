package mocks

import (
	"github.com/stretchr/testify/mock"
	"org.alex859/stockprices/domain"
)

type PricesSaverMock struct {
	mock.Mock
}

func (mock *PricesSaverMock) SavePrices(ticker domain.Ticker, prices *domain.PriceHistory) error {
	args := mock.Called(ticker, prices)
	return getErrorValueOrNil(args.Get(0))
}


