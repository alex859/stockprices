// Code generated by mockery v1.0.0. DO NOT EDIT.

package googlefinance

import mock "github.com/stretchr/testify/mock"

// MockPricesFetcher is an autogenerated mock type for the PricesFetcher type
type MockPricesFetcher struct {
	mock.Mock
}

// FetchPrices provides a mock function with given fields: market, symbol, period
func (_m *MockPricesFetcher) FetchPrices(market string, symbol string, period Period) (Response, error) {
	ret := _m.Called(market, symbol, period)

	var r0 Response
	if rf, ok := ret.Get(0).(func(string, string, Period) Response); ok {
		r0 = rf(market, symbol, period)
	} else {
		r0 = ret.Get(0).(Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, Period) error); ok {
		r1 = rf(market, symbol, period)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
