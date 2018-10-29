package pricesprovider

import "org.alex859/stockprices/domain"

type (
	// Retrieve prices data from Google Finance.
	GoogleFinancePricesFetcher interface {
		FetchPrices(market string, symbol string, period domain.Period) (domain.GoogleFinanceResponse, error)
	}

	GoogleFinanceResponseToPriceHistoryConverter interface {
		ConvertToPriceHistory(response domain.GoogleFinanceResponse) (domain.PriceHistory, error)
	}

	GoogleFinanceResponseToCurrentPriceConverter interface {
		ConvertToCurrentPrice(response domain.GoogleFinanceResponse) (domain.CurrentPrice, error)
	}

	GoogleFinanceResponseConverter interface {
		GoogleFinanceResponseToPriceHistoryConverter
		GoogleFinanceResponseToCurrentPriceConverter
	}
)