package googlefinance

import "org.alex859/stockprices/domain/entity"

type (
	// PricesFetcher goes off to GoogleFinance to retrieve prices data.
	PricesFetcher interface {
		FetchPrices(market string, symbol string, period Period) (Response, error)
	}

	// ResponseToPriceHistoryConverter converts the response from GoogleFinance into a PriceHistory.
	ResponseToPriceHistoryConverter interface {
		ConvertToPriceHistory(response Response) (entity.PriceHistory, error)
	}

	// ResponseToCurrentPriceConverter converts the response from GoogleFinance into a PriceHistory.
	ResponseToCurrentPriceConverter interface {
		ConvertToCurrentPrice(response Response) (entity.CurrentPrice, error)
	}

	// ResponseConverter allows to get current price and price history.
	ResponseConverter interface {
		ResponseToCurrentPriceConverter
		ResponseToPriceHistoryConverter
	}

	// Response models data coming from Google Finance.
	Response struct {
		Currency      string
		Name          string
		Symbol        string
		Market        string
		LastPrice     string
		LastPriceTime string
		PricesRows    []PriceRow
	}

	// PriceRow is a Single time/price.
	PriceRow struct {
		Price string
		Time  string
	}
)
