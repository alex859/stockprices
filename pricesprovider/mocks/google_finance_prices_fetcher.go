// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import domain "org.alex859/stockprices/domain"
import mock "github.com/stretchr/testify/mock"

// GoogleFinancePricesFetcher is an autogenerated mock type for the GoogleFinancePricesFetcher type
type GoogleFinancePricesFetcher struct {
	mock.Mock
}

// FetchPrices provides a mock function with given fields: market, symbol, period
func (_m *GoogleFinancePricesFetcher) FetchPrices(market string, symbol string, period domain.Period) (domain.GoogleFinanceResponse, error) {
	ret := _m.Called(market, symbol, period)

	var r0 domain.GoogleFinanceResponse
	if rf, ok := ret.Get(0).(func(string, string, domain.Period) domain.GoogleFinanceResponse); ok {
		r0 = rf(market, symbol, period)
	} else {
		r0 = ret.Get(0).(domain.GoogleFinanceResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, domain.Period) error); ok {
		r1 = rf(market, symbol, period)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}