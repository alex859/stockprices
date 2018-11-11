package usecase

import (
	"testing"
	"org.alex859/stockprices/domain/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	"org.alex859/stockprices/domain/usecase/mocks"
)

func Test_GetCurrentPrices_WHEN_NoTickers_THEN_ReturnEmptyMap(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	useCase := NewGetCurrentPricesUseCase(priceProvider, 1)
	result, err := useCase.GetCurrentPrices([]entity.Ticker{})

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]entity.CurrentPrice{}, result)
	}
}

func Test_GetCurrentPrices_WHEN_ErrorQueryingOneTicker_THEN_ReturnError(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	ticker1 := entity.Ticker{Symbol:"ANP", Market:"LON"}
	var providedHistory entity.CurrentPrice
	priceProvider.On("GetCurrentPrice", ticker1).Return(providedHistory, errors.New("an error occurred"))
	useCase := NewGetCurrentPricesUseCase(priceProvider, 1)
	_, err := useCase.GetCurrentPrices([]entity.Ticker{ticker1})

	assert.Error(t, err)
}

func Test_GetCurrentPrices_WHEN_OKQueryingOneTicker_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	ticker1 := entity.Ticker{Symbol:"ANP", Market:"LON"}

	priceProvider.On("GetCurrentPrice", ticker1).Return(entity.CurrentPrice{TickerInfo:tickerInfoAnp}, nil)
	useCase := NewGetCurrentPricesUseCase(priceProvider, 1)
	result, err := useCase.GetCurrentPrices([]entity.Ticker{ticker1})

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]entity.CurrentPrice{"LON:ANP": {TickerInfo: tickerInfoAnp}}, result)
	}

}

func Test_GetCurrentPrices_WHEN_OKQueryingMultipleTicker_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	ticker1 := entity.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := entity.Ticker{Symbol:"SDRY", Market:"LON"}
	priceProvider.On("GetCurrentPrice", ticker1).Return(entity.CurrentPrice{TickerInfo: tickerInfoAnp}, nil)
	priceProvider.On("GetCurrentPrice", ticker2).Return(entity.CurrentPrice{TickerInfo: tickerInfoSdry}, nil)
	useCase := NewGetCurrentPricesUseCase(priceProvider, 1)
	result, err := useCase.GetCurrentPrices([]entity.Ticker{ticker1, ticker2})

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]entity.CurrentPrice{
			"LON:ANP": {TickerInfo: tickerInfoAnp},
			"LON:SDRY": {TickerInfo: tickerInfoSdry},
		}, result)
	}

}

func Test_GetCurrentPrices_WHEN_OKQueryingMultipleTickerAndMultipleWorkers_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	ticker1 := entity.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := entity.Ticker{Symbol:"SDRY", Market:"LON"}
	priceProvider.On("GetCurrentPrice", ticker1).Return(entity.CurrentPrice{TickerInfo: tickerInfoAnp}, nil)
	priceProvider.On("GetCurrentPrice", ticker2).Return(entity.CurrentPrice{TickerInfo: tickerInfoSdry}, nil)
	useCase := NewGetCurrentPricesUseCase(priceProvider, 2)
	result, err := useCase.GetCurrentPrices([]entity.Ticker{ticker1, ticker2})

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]entity.CurrentPrice{
			"LON:ANP": {TickerInfo: tickerInfoAnp},
			"LON:SDRY": {TickerInfo: tickerInfoSdry},
		}, result)
	}

}

func Test_GetCurrentPrices_WHEN_ErrorQueryingAllMultipleTicker_THEN_ReturnError(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	ticker1 := entity.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := entity.Ticker{Symbol:"SDRY", Market:"LON"}
	var noResult entity.CurrentPrice
	priceProvider.On("GetCurrentPrice", ticker1).Return(noResult, errors.New("an error occurred"))
	priceProvider.On("GetCurrentPrice", ticker2).Return(noResult, errors.New("an error occurred"))
	useCase := NewGetCurrentPricesUseCase(priceProvider, 1)
	_, err := useCase.GetCurrentPrices([]entity.Ticker{ticker1, ticker2})

	assert.Error(t, err)
}

func Test_GetCurrentPrices_WHEN_ErrorQueryingOneMultipleTicker_THEN_ReturnOneCorrect(t *testing.T) {
	priceProvider := &mocks.CurrentPriceProvider{}
	ticker1 := entity.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := entity.Ticker{Symbol:"SDRY", Market:"LON"}
	var noResult entity.CurrentPrice
	priceProvider.On("GetCurrentPrice", ticker1).Return(entity.CurrentPrice{TickerInfo: tickerInfoAnp}, nil)
	priceProvider.On("GetCurrentPrice", ticker2).Return(noResult, errors.New("an error occurred"))
	useCase := NewGetCurrentPricesUseCase(priceProvider, 1)
	result, err := useCase.GetCurrentPrices([]entity.Ticker{ticker1, ticker2})

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]entity.CurrentPrice{
			"LON:ANP": {TickerInfo: tickerInfoAnp},
		}, result)
	}

}