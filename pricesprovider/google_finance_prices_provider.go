package pricesprovider

import (
	"org.alex859/stockprices/domain"
	"github.com/pkg/errors"
)

// Retrieves prices from Google Finance.
type googleFinanceHistoricalPricesProvider struct {
	googlePricesClient GoogleFinancePricesFetcher
	converter   GoogleFinanceResponseConverter
}

func NewGoogleFinancePricesProvider(client GoogleFinancePricesFetcher, converter GoogleFinanceResponseConverter) *googleFinanceHistoricalPricesProvider {
	return &googleFinanceHistoricalPricesProvider{googlePricesClient: client, converter: converter}
}

func (provider *googleFinanceHistoricalPricesProvider) GetHistoricalPrices(ticker domain.Ticker, interval domain.DateInterval) (result domain.PriceHistory, err error) {
	var googleResponse domain.GoogleFinanceResponse
	if googleResponse, err = provider.googlePricesClient.FetchPrices(ticker.Market, ticker.Symbol, interval.ToPeriod()); err == nil {
		if result, err = provider.converter.ConvertToPriceHistory(googleResponse); err == nil {
			result.Prices = result.Prices.FilterByInterval(interval)
			return result, nil
		}
	}

	return result, errors.Wrap(err, "unable to get historical prices")
}

func (provider *googleFinanceHistoricalPricesProvider) GetCurrentPrice(ticker domain.Ticker) (result domain.CurrentPrice, err error) {
	var googleResponse domain.GoogleFinanceResponse
	if googleResponse, err = provider.googlePricesClient.FetchPrices(ticker.Market, ticker.Symbol, domain.OneDay); err == nil {
		if result, err = provider.converter.ConvertToCurrentPrice(googleResponse); err == nil {
			return result, nil
		}
	}

	return result, errors.Wrap(err, "unable to get current price")
}
