package googlefinance

import (
	"reflect"
	"testing"

	"org.alex859/stockprices/domain"
	"time"
)

var goodResponse = domain.GoogleFinanceResponse{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"6 Sep, 15:04 BST",
	PricesRows:[]domain.GoogleFinancePriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongLastPrice = domain.GoogleFinanceResponse{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.ss25", LastPriceTime:"6 Sep, 14:04 BST",
	PricesRows:[]domain.GoogleFinancePriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongLastPriceTime = domain.GoogleFinanceResponse{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"6 Sep, YY BST",
	PricesRows:[]domain.GoogleFinancePriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongTime = domain.GoogleFinanceResponse{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"6 Sep, 15:04 BST",
	PricesRows:[]domain.GoogleFinancePriceRow{
		{Time:"2AA84030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongPrice = domain.GoogleFinanceResponse{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"6 Sep, 15:04 BST",
	PricesRows:[]domain.GoogleFinancePriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26sss"},
	},
}
var date1, _ = time.Parse("02-01-2006T15:04:05.000Z", "23-08-2018T16:30:00.000Z")
var date2, _ = time.Parse("02-01-2006T15:04:05.000Z", "24-08-2018T09:10:00.000Z")
var date3, _ = time.Parse("02-01-2006T15:04:05.000Z", "06-09-2018T14:04:00.000Z")
var goodPriceHistory = domain.PriceHistory{
	TickerInfo: domain.TickerInfo{Name:"Anpario", Currency:"GBX", Ticker:domain.Ticker{Market:"LON", Symbol:"ANP"}},
	Prices:domain.PriceList{
		{Price:20.25, Time: date1},
		{Price:20.26, Time: date2},
	},
}

var goodCurrentPrice = domain.CurrentPrice{
	TickerInfo: domain.TickerInfo{Name:"Anpario", Currency:"GBX", Ticker:domain.Ticker{Market:"LON", Symbol:"ANP"}},
	Time:       date3,
	Price:      12.25,
}

func Test_googleFinanceResponseConverter_ConvertToPriceHistory(t *testing.T) {
	type args struct {
		response domain.GoogleFinanceResponse
	}
	tests := []struct {
		name      string
		args      args
		want      domain.PriceHistory
		wantErr   bool
	}{
		{"All good", args{response:goodResponse}, goodPriceHistory, false},
		{"When wrong last price should return correct price history", args{response:wrongLastPrice}, goodPriceHistory, false},
		{"When wrong last time should return correct price history", args{response:wrongLastPriceTime}, goodPriceHistory, false},
		{"When wrong price should return error", args{response:wrongPrice}, domain.PriceHistory{}, true},
		{"When wrong time should return error", args{response:wrongTime}, domain.PriceHistory{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := &googleFinanceResponseConverter{}
			got, err := converter.ConvertToPriceHistory(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("googleFinanceResponseConverter.ConvertToPriceHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("googleFinanceResponseConverter.ConvertToPriceHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_googleFinanceResponseConverter_ConvertToCurrentPrice(t *testing.T) {
	now = func() time.Time {
		return date1
	}
	type args struct {
		response domain.GoogleFinanceResponse
	}
	tests := []struct {
		name      string
		args      args
		want      domain.CurrentPrice
		wantErr   bool
	}{
		{"All good", args{response:goodResponse}, goodCurrentPrice, false},
		{"When wrong last price should return error", args{response:wrongLastPrice}, domain.CurrentPrice{}, true},
		{"When wrong last time should return error", args{response:wrongLastPriceTime}, domain.CurrentPrice{}, true},
		{"When wrong price should return correct current price", args{response:wrongPrice}, goodCurrentPrice, false},
		{"When wrong time should return correct current price", args{response:wrongTime}, goodCurrentPrice, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := &googleFinanceResponseConverter{}
			got, err := converter.ConvertToCurrentPrice(tt.args.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("googleFinanceResponseConverter.ConvertToCurrentPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("googleFinanceResponseConverter.ConvertToCurrentPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
