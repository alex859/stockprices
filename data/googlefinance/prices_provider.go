package googlefinance

import (
	"org.alex859/stockprices/domain/entity"
	"github.com/pkg/errors"
)

// Retrieves prices from Google Finance.
type googleFinanceHistoricalPricesProvider struct {
	googlePricesClient PricesFetcher
	converter   ResponseConverter
}

// NewGoogleFinancePricesProvider Creates a new GoogleFinance prices provider.
func NewGoogleFinancePricesProvider(client PricesFetcher, converter ResponseConverter) *googleFinanceHistoricalPricesProvider {
	return &googleFinanceHistoricalPricesProvider{googlePricesClient: client, converter: converter}
}

func (provider *googleFinanceHistoricalPricesProvider) GetHistoricalPrices(ticker entity.Ticker, interval entity.DateInterval) (result entity.PriceHistory, err error) {
	var googleResponse Response
	if googleResponse, err = provider.googlePricesClient.FetchPrices(ticker.Market, ticker.Symbol, FromDateInterval(interval)); err == nil {
		if result, err = provider.converter.ConvertToPriceHistory(googleResponse); err == nil {
			result.Prices = result.Prices.FilterByInterval(interval)
			return result, nil
		}
	}

	return result, errors.Wrap(err, "unable to get historical prices")
}

func (provider *googleFinanceHistoricalPricesProvider) GetCurrentPrice(ticker entity.Ticker) (result entity.CurrentPrice, err error) {
	var googleResponse Response
	if googleResponse, err = provider.googlePricesClient.FetchPrices(ticker.Market, ticker.Symbol, OneDay); err == nil {
		if result, err = provider.converter.ConvertToCurrentPrice(googleResponse); err == nil {
			return result, nil
		}
	}

	return result, errors.Wrap(err, "unable to get current price")
}
