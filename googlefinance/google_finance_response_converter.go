package googlefinance

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"org.alex859/stockprices/domain"
)

// Default implementation of conversion logic from google finance response to PriceHistory and CurrentPrice.
type googleFinanceResponseConverter struct {
}

// for testing
var now = time.Now

var dateLayout = "Jan 2, 3:04 PM MST 2006"

// NewGoogleFinanceResponseConverter creates a new googleFinanceResponseConverter.
func NewGoogleFinanceResponseConverter() *googleFinanceResponseConverter {
	return &googleFinanceResponseConverter{}
}

func (converter *googleFinanceResponseConverter) ConvertToPriceHistory(response domain.GoogleFinanceResponse) (result domain.PriceHistory, err error) {
	var priceList domain.PriceList
	if priceList, err = convertPriceRows(response.PricesRows); err == nil {
		result = domain.PriceHistory{
			TickerInfo: convertTickerInfo(response),
			Prices:     priceList,
		}
	}
	err = errors.Wrap(err, "error converting to price history")
	return
}

func (converter *googleFinanceResponseConverter) ConvertToCurrentPrice(response domain.GoogleFinanceResponse) (result domain.CurrentPrice, err error) {
	var price float64
	if price, err = convertPrice(response.LastPrice); err == nil {
		var lastTime time.Time
		currentYear := now().Year()
		if lastTime, err = time.Parse(dateLayout, fmt.Sprintf("%s %v", response.LastPriceTime, currentYear)); err == nil {
			result = domain.CurrentPrice{
				TickerInfo: convertTickerInfo(response),
				Price:      price,
				Time:       lastTime.In(time.UTC),
			}
		}
	}
	return result, errors.Wrap(err, "error converting to current price")
}

func convertPriceRows(priceRows []domain.GoogleFinancePriceRow) (domain.PriceList, error) {
	result := make(domain.PriceList, len(priceRows))
	for i, p := range priceRows {
		if pp, err := convertPriceRow(p); err == nil {
			result[i] = pp
		} else {
			return result, errors.Wrap(err, "error converting price rows")
		}

	}
	return result, nil
}

func convertPriceRow(priceRow domain.GoogleFinancePriceRow) (result domain.PricePoint, err error) {
	var p float64
	if p, err = convertPrice(priceRow.Price); err == nil {
		var t float64
		if t, err = strconv.ParseFloat(priceRow.Time, 64); err == nil {
			result = domain.PricePoint{
				Price: p,
				Time:  time.Unix(60*(int64(t)), 0).In(time.UTC),
			}
		}
	}
	return result, errors.Wrap(err, "error converting price row")
}

func convertTickerInfo(response domain.GoogleFinanceResponse) domain.TickerInfo {
	return domain.TickerInfo{
		Ticker:   convertTicker(response),
		Currency: response.Currency,
		Name:     response.Name,
	}
}

func convertTicker(response domain.GoogleFinanceResponse) domain.Ticker {
	return domain.Ticker{
		Symbol: response.Symbol,
		Market: response.Market,
	}
}

func convertPrice(data string) (float64, error) {
	return strconv.ParseFloat(removeCommas(data), 64)
}

func removeCommas(str interface{}) string {
	return strings.Replace(str.(string), ",", "", -1)
}