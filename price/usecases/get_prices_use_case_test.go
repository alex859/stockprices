package usecases

import (
	"testing"
	"org.alex859/stockprices/price/usecases/mocks"
	"org.alex859/stockprices/domain"
	"errors"
	"github.com/stretchr/testify/assert"
	"time"
)

var from, _ = time.Parse("2006-01-02", "2018-01-01")
var to, _ = time.Parse("2006-01-02", "2018-01-10")
func Test_FetchPrices_WHEN_NoTickers_THEN_ReturnEmptyMap(t *testing.T) {
	priceProvider := &mocks.PriceProviderMock{}
	useCase := NewGetPricesUseCase(priceProvider)
	result, err := useCase.GetMultiplePrices([]domain.Ticker{}, from, to)

	assert.Equal(t, map[domain.Ticker]*domain.PriceHistory{}, result)
	assert.Nil(t, err)
}

func Test_FetchPrices_WHEN_ErrorQueryingOneTicker_THEN_ReturnError(t *testing.T) {
	priceProvider := &mocks.PriceProviderMock{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	priceProvider.On("FetchPrices", ticker1, from, to).Return(nil, errors.New("an error occurred"))
	useCase := NewGetPricesUseCase(priceProvider)
	result, err := useCase.GetMultiplePrices([]domain.Ticker{ticker1}, from, to)

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func Test_FetchPrices_WHEN_OKQueryingOneTicker_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.PriceProviderMock{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	priceProvider.On("FetchPrices", ticker1, from, to).Return(&domain.PriceHistory{Ticker:ticker1}, nil)
	useCase := NewGetPricesUseCase(priceProvider)
	result, err := useCase.GetMultiplePrices([]domain.Ticker{ticker1}, from, to)

	assert.Equal(t, map[domain.Ticker]*domain.PriceHistory{ticker1: {Ticker: ticker1}}, result)
	assert.Nil(t, err)
}

func Test_FetchPrices_WHEN_OKQueryingMultipleTicker_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.PriceProviderMock{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	priceProvider.On("FetchPrices", ticker1, from, to).Return(&domain.PriceHistory{Ticker:ticker1}, nil)
	priceProvider.On("FetchPrices", ticker2, from, to).Return(&domain.PriceHistory{Ticker:ticker2}, nil)
	useCase := NewGetPricesUseCase(priceProvider)
	result, err := useCase.GetMultiplePrices([]domain.Ticker{ticker1, ticker2}, from, to)

	assert.Nil(t, err)
	assert.Equal(t, map[domain.Ticker]*domain.PriceHistory{
		ticker1: {Ticker: ticker1},
		ticker2: {Ticker: ticker2},
	}, result)
}

func Test_FetchPrices_WHEN_ErrorQueryingAllMultipleTicker_THEN_ReturnError(t *testing.T) {
	priceProvider := &mocks.PriceProviderMock{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	priceProvider.On("FetchPrices", ticker1, from, to).Return(nil, errors.New("an error occurred"))
	priceProvider.On("FetchPrices", ticker2, from, to).Return(nil, errors.New("an error occurred"))
	useCase := NewGetPricesUseCase(priceProvider)
	result, err := useCase.GetMultiplePrices([]domain.Ticker{ticker1, ticker2}, from, to)

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func Test_FetchPrices_WHEN_ErrorQueryingOneMultipleTicker_THEN_ReturnOneCorrect(t *testing.T) {
	priceProvider := &mocks.PriceProviderMock{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	priceProvider.On("FetchPrices", ticker1, from, to).Return(&domain.PriceHistory{Ticker:ticker1}, nil)
	priceProvider.On("FetchPrices", ticker2, from, to).Return(nil, errors.New("an error occurred"))
	useCase := NewGetPricesUseCase(priceProvider)
	result, err := useCase.GetMultiplePrices([]domain.Ticker{ticker1, ticker2}, from, to)

	assert.Equal(t, map[domain.Ticker]*domain.PriceHistory{
		ticker1: {Ticker: ticker1},
	}, result)
	assert.Nil(t, err)
}