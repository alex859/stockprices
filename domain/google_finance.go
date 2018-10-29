package domain

type (
	// Price data coming from Google Finance.
	GoogleFinanceResponse struct {
		Currency      string
		Name          string
		Symbol        string
		Market        string
		LastPrice     string
		LastPriceTime string
		PricesRows    []GoogleFinancePriceRow
	}

	// Single time/price.
	GoogleFinancePriceRow struct {
		Price string
		Time  string
	}
)
