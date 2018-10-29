package usecases

import (
	"testing"
	"org.alex859/stockprices/domain"
	"errors"
	"github.com/stretchr/testify/assert"
	"time"
	"org.alex859/stockprices/usecases/mocks"
)

var from, _ = time.Parse("2006-01-02", "2018-01-01")
var to, _ = time.Parse("2006-01-02", "2018-01-10")

var tickerInfoAnp = domain.TickerInfo{Ticker:domain.Ticker{Symbol:"ANP", Market:"LON"}, Currency:"GBX"}
var tickerInfoSdry = domain.TickerInfo{Ticker:domain.Ticker{Symbol:"SDRY", Market:"LON"}, Currency:"GBX"}

func Test_GetHistoricalPrices_WHEN_NoTickers_THEN_ReturnEmptyMap(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 1)
	interval, err := domain.NewDateInterval(from, to)
	result, err := useCase.GetHistoricalPrices([]domain.Ticker{}, interval)

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]domain.PriceHistory{}, result)
	}
}

func Test_GetHistoricalPrices_WHEN_ErrorQueryingOneTicker_THEN_ReturnError(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	interval, err := domain.NewDateInterval(from, to)
	var providedHistory domain.PriceHistory
	priceProvider.On("GetHistoricalPrices", ticker1, interval).Return(providedHistory, errors.New("an error occurred"))
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 1)
	_, err = useCase.GetHistoricalPrices([]domain.Ticker{ticker1}, interval)

	assert.Error(t, err)
}

func Test_GetHistoricalPrices_WHEN_OKQueryingOneTicker_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	interval, err := domain.NewDateInterval(from, to)

	priceProvider.On("GetHistoricalPrices", ticker1, interval).Return(domain.PriceHistory{TickerInfo:tickerInfoAnp}, nil)
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 1)
	result, err := useCase.GetHistoricalPrices([]domain.Ticker{ticker1}, interval)

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]domain.PriceHistory{"LON:ANP": {TickerInfo: tickerInfoAnp}}, result)
	}

}

func Test_GetHistoricalPrices_WHEN_OKQueryingMultipleTicker_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	interval, err := domain.NewDateInterval(from, to)
	priceProvider.On("GetHistoricalPrices", ticker1, interval).Return(domain.PriceHistory{TickerInfo: tickerInfoAnp}, nil)
	priceProvider.On("GetHistoricalPrices", ticker2, interval).Return(domain.PriceHistory{TickerInfo: tickerInfoSdry}, nil)
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 1)
	result, err := useCase.GetHistoricalPrices([]domain.Ticker{ticker1, ticker2}, interval)

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]domain.PriceHistory{
			"LON:ANP": {TickerInfo: tickerInfoAnp},
			"LON:SDRY": {TickerInfo: tickerInfoSdry},
		}, result)
	}

}

func Test_GetHistoricalPrices_WHEN_OKQueryingMultipleTickerAndMultipleWorkers_THEN_ReturnCorrect(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	interval, err := domain.NewDateInterval(from, to)
	priceProvider.On("GetHistoricalPrices", ticker1, interval).Return(domain.PriceHistory{TickerInfo: tickerInfoAnp}, nil)
	priceProvider.On("GetHistoricalPrices", ticker2, interval).Return(domain.PriceHistory{TickerInfo: tickerInfoSdry}, nil)
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 2)
	result, err := useCase.GetHistoricalPrices([]domain.Ticker{ticker1, ticker2}, interval)

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]domain.PriceHistory{
			"LON:ANP": {TickerInfo: tickerInfoAnp},
			"LON:SDRY": {TickerInfo: tickerInfoSdry},
		}, result)
	}

}

func Test_GetHistoricalPrices_WHEN_ErrorQueryingAllMultipleTicker_THEN_ReturnError(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	interval, err := domain.NewDateInterval(from, to)
	var noResult domain.PriceHistory
	priceProvider.On("GetHistoricalPrices", ticker1, interval).Return(noResult, errors.New("an error occurred"))
	priceProvider.On("GetHistoricalPrices", ticker2, interval).Return(noResult, errors.New("an error occurred"))
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 1)
	_, err = useCase.GetHistoricalPrices([]domain.Ticker{ticker1, ticker2}, interval)

	assert.Error(t, err)
}

func Test_GetHistoricalPrices_WHEN_ErrorQueryingOneMultipleTicker_THEN_ReturnOneCorrect(t *testing.T) {
	priceProvider := &mocks.HistoricalPricesProvider{}
	ticker1 := domain.Ticker{Symbol:"ANP", Market:"LON"}
	ticker2 := domain.Ticker{Symbol:"SDRY", Market:"LON"}
	interval, err := domain.NewDateInterval(from, to)
	var noResult domain.PriceHistory
	priceProvider.On("GetHistoricalPrices", ticker1, interval).Return(domain.PriceHistory{TickerInfo: tickerInfoAnp}, nil)
	priceProvider.On("GetHistoricalPrices", ticker2, interval).Return(noResult, errors.New("an error occurred"))
	useCase := NewGetHistoricalPricesUseCase(priceProvider, 1)
	result, err := useCase.GetHistoricalPrices([]domain.Ticker{ticker1, ticker2}, interval)

	if assert.NoError(t, err) {
		assert.Equal(t, map[string]domain.PriceHistory{
			"LON:ANP": {TickerInfo: tickerInfoAnp},
		}, result)
	}

}