package interfaces

import (
	"testing"
	"org.alex859/stockprices/price/usecases/mocks"
	"org.alex859/stockprices/domain"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"errors"
)

func Test_WhenAllFoundInCache_ShouldNotGoOffAndNotSave(t *testing.T) {
	ticker := domain.Ticker{Market: "LON", Symbol: "ANP"}
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var may6, _ = time.Parse("02-01-2006", "06-05-2018")
	var june6, _ = time.Parse("02-01-2006", "06-06-2018")
	var july6, _ = time.Parse("02-01-2006", "06-07-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	ppMay5 := domain.PricePoint{Price: 123, Time: may5}
	ppMay6 := domain.PricePoint{Price: 123, Time: may6}
	ppJune6 := domain.PricePoint{Price: 123, Time: june6}
	ppJuly6 := domain.PricePoint{Price: 123, Time: july6}
	ppJuly16 := domain.PricePoint{Price: 123, Time: july16}
	cachedResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 125, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppMay5, ppMay6, ppJune6, ppJuly6, ppJuly16},
	}
	priceProvider := &mocks.PriceProviderMock{}
	cachedPriceProvider := &mocks.PriceProviderMock{}
	cachedPriceProvider.On("FetchPrices", ticker, may5, july16).Return(&cachedResult, nil)
	priceSaver := &mocks.PricesSaverMock{}
	providerUT := NewCachingPricesProvider(priceProvider, cachedPriceProvider, priceSaver)

	result, err := providerUT.FetchPrices(ticker, may5, july16)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedResult, *result)
	priceProvider.AssertNotCalled(t, "FetchPrices", mock.Anything, mock.Anything)
}

func Test_WhenNothingFoundInCacheAndErrorGoingOff_ShouldReturnError(t *testing.T) {
	ticker := domain.Ticker{Market: "LON", Symbol: "ANP"}
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	priceProvider := &mocks.PriceProviderMock{}
	priceProvider.On("FetchPrices", ticker, may5, july16).Return(nil, errors.New("unable to access external service"))
	cachedPriceProvider := &mocks.PriceProviderMock{}
	cachedPriceProvider.On("FetchPrices", ticker, may5, july16).Return(nil, errors.New("unable to access cache"))
	priceSaver := &mocks.PricesSaverMock{}
	providerUT := NewCachingPricesProvider(priceProvider, cachedPriceProvider, priceSaver)

	result, err := providerUT.FetchPrices(ticker, may5, july16)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	priceProvider.AssertCalled(t, "FetchPrices", may5, july16)
	priceSaver.AssertNotCalled(t, "SavePrices", mock.Anything, mock.Anything)
}

func Test_WhenCacheAccessFails_ShouldGoOffAndSave(t *testing.T) {
	ticker := domain.Ticker{Market: "LON", Symbol: "ANP"}
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var may6, _ = time.Parse("02-01-2006", "06-05-2018")
	var june6, _ = time.Parse("02-01-2006", "06-06-2018")
	var july6, _ = time.Parse("02-01-2006", "06-07-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	ppMay5 := domain.PricePoint{Price: 123, Time: may5}
	ppMay6 := domain.PricePoint{Price: 123, Time: may6}
	ppJune6 := domain.PricePoint{Price: 123, Time: june6}
	ppJuly6 := domain.PricePoint{Price: 123, Time: july6}
	ppJuly16 := domain.PricePoint{Price: 123, Time: july16}
	priceProvider := &mocks.PriceProviderMock{}
	providerResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 125, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppMay5, ppMay6, ppJune6, ppJuly6, ppJuly16},
	}
	cachedPriceProvider := &mocks.PriceProviderMock{}
	cachedPriceProvider.On("FetchPrices", ticker, may5, july16).Return(nil, errors.New("unable to access the cache"))
	priceSaver := &mocks.PricesSaverMock{}
	providerUT := NewCachingPricesProvider(priceProvider, cachedPriceProvider, priceSaver)

	result, err := providerUT.FetchPrices(ticker, may5, july16)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, providerResult, *result)
	priceProvider.AssertCalled(t, "FetchPrices", may5, july16)
	priceSaver.AssertCalled(t, "SavePrices", ticker, &result)
}

func Test_WhenSomethingFoundInCache_MissingPricesAfter_ShouldGoOffMergeAndSave(t *testing.T) {
	ticker := domain.Ticker{Market: "LON", Symbol: "ANP"}
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var may6, _ = time.Parse("02-01-2006", "06-05-2018")
	var june6, _ = time.Parse("02-01-2006", "06-06-2018")
	var june7, _ = time.Parse("02-01-2006", "07-06-2018")
	var july6, _ = time.Parse("02-01-2006", "06-07-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	ppMay5 := domain.PricePoint{Price: 123, Time: may5}
	ppMay6 := domain.PricePoint{Price: 123, Time: may6}
	ppJune6 := domain.PricePoint{Price: 123, Time: june6}
	ppJuly6 := domain.PricePoint{Price: 123, Time: july6}
	ppJuly16 := domain.PricePoint{Price: 123, Time: july16}

	cachedPriceProvider := &mocks.PriceProviderMock{}
	priceProvider := &mocks.PriceProviderMock{}
	priceSaver := &mocks.PricesSaverMock{}

	providerResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 115, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppJuly6, ppJuly16},
	}
	priceProvider.On("FetchPrices", ticker, june7, july16).Return(&providerResult, nil)

	cachedResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 125, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppMay5, ppMay6, ppJune6},
	}
	cachedPriceProvider.On("FetchPrices", ticker, may5, july16).Return(&cachedResult, nil)

	providerUT := NewCachingPricesProvider(priceProvider, cachedPriceProvider, priceSaver)

	expectedResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 115, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppMay5, ppMay6, ppJune6, ppJuly6, ppJuly16},
	}
	result, err := providerUT.FetchPrices(ticker, may5, july16)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult, *result)
	priceProvider.AssertCalled(t, "FetchPrices", june7, july16)
	priceSaver.AssertCalled(t, "SavePrices", ticker, &expectedResult)
}

func Test_WhenSomethingFoundInCache_MissingPricesBefore_ShouldGoOffMergeAndSave(t *testing.T) {
	ticker := domain.Ticker{Market: "LON", Symbol: "ANP"}
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var may6, _ = time.Parse("02-01-2006", "06-05-2018")
	var june6, _ = time.Parse("02-01-2006", "06-06-2018")
	var july5, _ = time.Parse("02-01-2006", "05-07-2018")
	var july6, _ = time.Parse("02-01-2006", "06-07-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	ppMay5 := domain.PricePoint{Price: 123, Time: may5}
	ppMay6 := domain.PricePoint{Price: 123, Time: may6}
	ppJune6 := domain.PricePoint{Price: 123, Time: june6}
	ppJuly6 := domain.PricePoint{Price: 123, Time: july6}
	ppJuly16 := domain.PricePoint{Price: 123, Time: july16}

	cachedPriceProvider := &mocks.PriceProviderMock{}
	priceProvider := &mocks.PriceProviderMock{}
	priceSaver := &mocks.PricesSaverMock{}

	cachedResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 115, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppJuly6, ppJuly16},
	}
	cachedPriceProvider.On("FetchPrices", ticker, may5, july16).Return(&cachedResult, nil)

	providerResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 115, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppMay5, ppMay6, ppJune6},
	}
	priceProvider.On("FetchPrices", ticker, may5, july5).Return(&providerResult, nil)

	providerUT := NewCachingPricesProvider(priceProvider, cachedPriceProvider, priceSaver)

	expectedResult := domain.PriceHistory{
		Ticker: ticker, Currency: "GBX", LastPrice: 115, LastPriceTime: "last time", Name: "Anpario",
		Prices: domain.PriceList{ppMay5, ppMay6, ppJune6, ppJuly6, ppJuly16},
	}
	result, err := providerUT.FetchPrices(ticker, may5, july16)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResult, *result)
	priceProvider.AssertCalled(t, "FetchPrices",  may5, july5)
	priceSaver.AssertCalled(t, "SavePrices", ticker, &expectedResult)
}
