package googlefinance

import (
	"testing"
	"time"
	"org.alex859/stockprices/domain/entity"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetHistoricalPrices_WHEN_ErrorFromGoogle_THEN_ShouldReturnError(t *testing.T) {
	client := &MockPricesFetcher{}
	converter := &MockResponseConverter{}

	ticker := entity.Ticker{Symbol: "ANP", Market: "LON"}
	may5, _ := time.Parse("02-01-2006", "05-05-2018")
	may9, _ := time.Parse("02-01-2006", "09-05-2018")
	dateInterval, _ := entity.NewDateInterval(may5, may9)
	period := FromDateInterval(dateInterval)

	noResponse := Response{}
	client.On("FetchPrices", "LON", "ANP", period).Return(noResponse, errors.New("can't talk to google"))

	priceProvider := NewGoogleFinancePricesProvider(client, converter)
	_, err := priceProvider.GetHistoricalPrices(ticker, dateInterval)

	assert.Error(t, err)
	converter.AssertNotCalled(t, mock.Anything)
}

func Test_GetHistoricalPrices_WHEN_ErrorConverting_THEN_ShouldReturnError(t *testing.T) {
	client := &MockPricesFetcher{}
	converter := &MockResponseConverter{}

	ticker := entity.Ticker{Symbol: "ANP", Market: "LON"}
	may5, _ := time.Parse("02-01-2006", "05-05-2018")
	may9, _ := time.Parse("02-01-2006", "09-05-2018")
	dateInterval, _ := entity.NewDateInterval(may5, may9)
	period := FromDateInterval(dateInterval)

	malformedResponse := Response{LastPrice: "AAA"}
	client.On("FetchPrices", "LON", "ANP", period).Return(malformedResponse, nil)
	noConverted := entity.PriceHistory{}
	converter.On("ConvertToPriceHistory", malformedResponse).Return(noConverted, errors.New("malformed input"))

	priceProvider := NewGoogleFinancePricesProvider(client, converter)
	_, err := priceProvider.GetHistoricalPrices(ticker, dateInterval)

	assert.Error(t, err)
}

func Test_GetHistoricalPrices_WHEN_AllGood_THEN_ShouldReturnResult(t *testing.T) {
	client := &MockPricesFetcher{}
	converter := &MockResponseConverter{}

	ticker := entity.Ticker{Symbol: "ANP", Market: "LON"}
	may5, _ := time.Parse("02-01-2006", "05-05-2018")
	may9, _ := time.Parse("02-01-2006", "09-05-2018")
	may10, _ := time.Parse("02-01-2006", "10-05-2018")
	dateInterval, _ := entity.NewDateInterval(may5, may9)
	period := FromDateInterval(dateInterval)

	goodResponse := Response{LastPrice: "12.5"}
	client.On("FetchPrices", "LON", "ANP", period).Return(goodResponse, nil)
	converted := entity.PriceHistory{
		TickerInfo: entity.TickerInfo{Name: "Anpario", Ticker: ticker, Currency: ""},
		Prices: entity.PriceList{
			{Time: may5, Price: 12.25},
			{Time: may9, Price: 12.5},
			{Time: may10, Price: 11.25},
		},
	}
	converter.On("ConvertToPriceHistory", goodResponse).Return(converted, nil)

	priceProvider := NewGoogleFinancePricesProvider(client, converter)

	expected :=  entity.PriceHistory{
		TickerInfo: entity.TickerInfo{Name: "Anpario", Ticker: ticker, Currency: ""},
		Prices: entity.PriceList{
			{Time: may5, Price: 12.25},
			{Time: may9, Price: 12.5},
		},
	}
	if result, err := priceProvider.GetHistoricalPrices(ticker, dateInterval); assert.NoError(t, err) {
		assert.Equal(t, expected, result)
	}
}

func Test_GetCurrentPrice_WHEN_ErrorFromGoogle_THEN_ShouldReturnError(t *testing.T) {
	client := &MockPricesFetcher{}
	converter := &MockResponseConverter{}

	ticker := entity.Ticker{Symbol: "ANP", Market: "LON"}

	noResponse := Response{}
	client.On("FetchPrices", "LON", "ANP", mock.Anything).Return(noResponse, errors.New("can't talk to google"))

	priceProvider := NewGoogleFinancePricesProvider(client, converter)
	_, err := priceProvider.GetCurrentPrice(ticker)

	assert.Error(t, err)
	converter.AssertNotCalled(t, mock.Anything)
}

func Test_GetCurrentPrice_WHEN_ErrorConverting_THEN_ShouldReturnError(t *testing.T) {
	client := &MockPricesFetcher{}
	converter := &MockResponseConverter{}

	ticker := entity.Ticker{Symbol: "ANP", Market: "LON"}

	malformedResponse := Response{LastPrice: "AAA"}
	client.On("FetchPrices", "LON", "ANP", mock.Anything).Return(malformedResponse, nil)
	noConverted := entity.CurrentPrice{}
	converter.On("ConvertToCurrentPrice", malformedResponse).Return(noConverted, errors.New("malformed input"))

	priceProvider := NewGoogleFinancePricesProvider(client, converter)
	_, err := priceProvider.GetCurrentPrice(ticker)

	assert.Error(t, err)
}

func Test_GetCurrentPrice_WHEN_AllGood_THEN_ShouldReturnResult(t *testing.T) {
	client := &MockPricesFetcher{}
	converter := &MockResponseConverter{}

	ticker := entity.Ticker{Symbol: "ANP", Market: "LON"}

	goodResponse := Response{LastPrice: "12.5"}
	client.On("FetchPrices", "LON", "ANP", mock.Anything).Return(goodResponse, nil)
	converted := entity.CurrentPrice{
		TickerInfo: entity.TickerInfo{Name: "Anpario", Ticker: ticker, Currency: ""},
	}
	converter.On("ConvertToCurrentPrice", goodResponse).Return(converted, nil)

	priceProvider := NewGoogleFinancePricesProvider(client, converter)

	if result, err := priceProvider.GetCurrentPrice(ticker); assert.NoError(t, err) {
		assert.Equal(t, converted, result)
	}
}
