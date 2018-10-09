package domain

import (
	"time"
	"sort"
	"errors"
	"fmt"
)

type PriceHistory struct {
	Name          string    `json:"name"`
	Ticker        Ticker    `json:"ticker"`
	Currency      string    `json:"currency"`
	LastPrice     float64   `json:"current_price"`
	LastPriceTime string    `json:"last_price_time"`
	Prices        PriceList `json:"prices"`
}

// Here we call LON:SDRY Ticker where LON is Market and SDRY is the Symbol
type Ticker struct {
	Market string `json:"market"`
	Symbol string `json:"symbol"`
}
func (t *Ticker) String() string {
	return fmt.Sprintf("%s:%s", t.Market, t.Symbol)
}

type PricePoint struct {
	Price float64   `json:"price"`
	Time  time.Time `json:"time"`
}

type PriceList []PricePoint

func (pl PriceList) On(day time.Time) (PricePoint, error) {
	if len(pl) == 0 {
		return PricePoint{}, errors.New("unable to get price on date for empty price history")
	}

	sort.Slice(pl, func(i, j int) bool {
		return pl[i].Time.Before(pl[j].Time)
	})

	dayToSearch := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, time.UTC)
	var i int
	for i = 0; i < len(pl) && pl[i].Time.Before(dayToSearch); i++ {

	}

	if i == 0 {
		return PricePoint{}, errors.New("nothing found")
	}

	// we should check if the day param is too far ahead
	return pl[i-1], nil
}
