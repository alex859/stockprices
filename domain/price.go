package domain

import (
	"time"
	"fmt"
)

type (
	// Here we call LON:SDRY Ticker where LON is Market and SDRY is the Symbol
	Ticker struct {
		Market string `json:"market"`
		Symbol string `json:"symbol"`
	}

	TickerInfo struct {
		Name     string `json:"name"`
		Ticker   Ticker `json:"ticker"`
		Currency string `json:"currency"`
	}

	CurrentPrice struct {
		TickerInfo
		Price float64   `json:"price"`
		Time  time.Time `json:"time"`
	}

	PriceHistory struct {
		TickerInfo
		Prices PriceList `json:"prices"`
	}
)

func (t Ticker) String() string {
	return fmt.Sprintf("%s:%s", t.Market, t.Symbol)
}
