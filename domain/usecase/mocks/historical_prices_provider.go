// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entity "org.alex859/stockprices/domain/entity"
import mock "github.com/stretchr/testify/mock"

// HistoricalPricesProvider is an autogenerated mock type for the HistoricalPricesProvider type
type HistoricalPricesProvider struct {
	mock.Mock
}

// GetHistoricalPrices provides a mock function with given fields: ticker, dateInterval
func (_m *HistoricalPricesProvider) GetHistoricalPrices(ticker entity.Ticker, dateInterval entity.DateInterval) (entity.PriceHistory, error) {
	ret := _m.Called(ticker, dateInterval)

	var r0 entity.PriceHistory
	if rf, ok := ret.Get(0).(func(entity.Ticker, entity.DateInterval) entity.PriceHistory); ok {
		r0 = rf(ticker, dateInterval)
	} else {
		r0 = ret.Get(0).(entity.PriceHistory)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.Ticker, entity.DateInterval) error); ok {
		r1 = rf(ticker, dateInterval)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}