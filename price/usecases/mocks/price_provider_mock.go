package mocks

import (
	"github.com/stretchr/testify/mock"
	"org.alex859/stockprices/domain"
	"time"
)

type PriceProviderMock struct {
	mock.Mock
}

func (mock *PriceProviderMock) FetchPrices(ticker domain.Ticker, from time.Time, to time.Time) (*domain.PriceHistory, error) {
	args := mock.Called(ticker, from, to)
	return getPriceHistoryOrNil(args.Get(0)), getErrorValueOrNil(args.Get(1))
}

func getPriceHistoryOrNil(v interface{}) *domain.PriceHistory{
	if v == nil {
		return nil
	}

	return v.(*domain.PriceHistory)
}

func getErrorValueOrNil(v interface{}) error{
	if v == nil {
		return nil
	}

	return v.(error)
}
