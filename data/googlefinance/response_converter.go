package googlefinance

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"org.alex859/stockprices/domain/entity"
)

// Default implementation of conversion logic from google finance response to PriceHistory and CurrentPrice.
type googleFinanceResponseConverter struct {
}

// for testing
var now = time.Now

var dateLayouts = []string{
	"Jan 2, 3:04 PM MST 2006",
	"2 Jan, 15:04 GMT 2006",
}

// NewGoogleFinanceResponseConverter creates a new googleFinanceResponseConverter.
func NewGoogleFinanceResponseConverter() *googleFinanceResponseConverter {
	return &googleFinanceResponseConverter{}
}

func (converter *googleFinanceResponseConverter) ConvertToPriceHistory(response Response) (result entity.PriceHistory, err error) {
	var priceList entity.PriceList
	if priceList, err = convertPriceRows(response.PricesRows); err == nil {
		result = entity.PriceHistory{
			TickerInfo: convertTickerInfo(response),
			Prices:     priceList,
		}
	}
	err = errors.Wrap(err, "error converting to price history")
	return
}

func (converter *googleFinanceResponseConverter) ConvertToCurrentPrice(response Response) (result entity.CurrentPrice, err error) {
	var price float64
	if price, err = convertPrice(response.LastPrice); err == nil {
		var lastTime time.Time
		currentYear := now().Year()
		if lastTime, err = readTime(fmt.Sprintf("%s %v", response.LastPriceTime, currentYear)); err == nil {
			result = entity.CurrentPrice{
				TickerInfo: convertTickerInfo(response),
				Price:      price,
				Time:       lastTime.In(time.UTC),
			}
		}
	}
	return result, errors.Wrap(err, "error converting to current price")
}

func readTime(str string) (result time.Time, err error) {
	for _, layout := range dateLayouts {
		result, err = time.Parse(layout, str)
		if err == nil {
			return result, err
		}
	}

	return result, err
}

func convertPriceRows(priceRows []PriceRow) (entity.PriceList, error) {
	result := make(entity.PriceList, len(priceRows))
	for i, p := range priceRows {
		if pp, err := convertPriceRow(p); err == nil {
			result[i] = pp
		} else {
			return result, errors.Wrap(err, "error converting price rows")
		}

	}
	return result, nil
}

func convertPriceRow(priceRow PriceRow) (result entity.PricePoint, err error) {
	var p float64
	if p, err = convertPrice(priceRow.Price); err == nil {
		var t float64
		if t, err = strconv.ParseFloat(priceRow.Time, 64); err == nil {
			result = entity.PricePoint{
				Price: p,
				Time:  time.Unix(60*(int64(t)), 0).In(time.UTC),
			}
		}
	}
	return result, errors.Wrap(err, "error converting price row")
}

func convertTickerInfo(response Response) entity.TickerInfo {
	return entity.TickerInfo{
		Ticker:   convertTicker(response),
		Currency: response.Currency,
		Name:     response.Name,
	}
}

func convertTicker(response Response) entity.Ticker {
	return entity.Ticker{
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