package googlefinance

import (
	"org.alex859/stockprices/domain/entity"
	"reflect"
	"testing"
	"time"
)

var goodResponse = Response{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"Sep 6, 2:04 PM GMT",
	PricesRows:[]PriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var goodResponseAnotherDateFormat = Response{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"6 Sep, 14:04 GMT",
	PricesRows:[]PriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongLastPrice = Response{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.ss25", LastPriceTime:"6 Sep, 14:04 BST",
	PricesRows:[]PriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongLastPriceTime = Response{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"6 Sep, YY BST",
	PricesRows:[]PriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongTime = Response{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"Sep 6, 2:04 PM GMT",
	PricesRows:[]PriceRow{
		{Time:"2AA84030", Price:"20.25"},
		{Time:"25585030", Price:"20.26"},
	},
}

var wrongPrice = Response{
	Symbol:"ANP", Market:"LON", Currency:"GBX", Name:"Anpario", LastPrice:"12.25", LastPriceTime:"Sep 6, 2:04 PM GMT",
	PricesRows:[]PriceRow{
		{Time:"25584030", Price:"20.25"},
		{Time:"25585030", Price:"20.26sss"},
	},
}
var date1, _ = time.Parse("02-01-2006T15:04:05.000Z", "23-08-2018T16:30:00.000Z")
var date2, _ = time.Parse("02-01-2006T15:04:05.000Z", "24-08-2018T09:10:00.000Z")
var date3, _ = time.Parse("02-01-2006T15:04:05.000Z", "06-09-2018T14:04:00.000Z")
var goodPriceHistory = entity.PriceHistory{
	TickerInfo: entity.TickerInfo{Name:"Anpario", Currency:"GBX", Ticker:entity.Ticker{Market:"LON", Symbol:"ANP"}},
	Prices:entity.PriceList{
		{Price:20.25, Time: date1},
		{Price:20.26, Time: date2},
	},
}

var goodCurrentPrice = entity.CurrentPrice{
	TickerInfo: entity.TickerInfo{Name:"Anpario", Currency:"GBX", Ticker:entity.Ticker{Market:"LON", Symbol:"ANP"}},
	Time:       date3,
	Price:      12.25,
}

func Test_googleFinanceResponseConverter_ConvertToPriceHistory(t *testing.T) {
	type args struct {
		response Response
	}
	tests := []struct {
		name      string
		args      args
		want      entity.PriceHistory
		wantErr   bool
	}{
		{"All good", args{response:goodResponse}, goodPriceHistory, false},
		{"When wrong last price should return correct price history", args{response:wrongLastPrice}, goodPriceHistory, false},
		{"When wrong last time should return correct price history", args{response:wrongLastPriceTime}, goodPriceHistory, false},
		{"When wrong price should return error", args{response:wrongPrice}, entity.PriceHistory{}, true},
		{"When wrong time should return error", args{response:wrongTime}, entity.PriceHistory{}, true},
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
		response Response
	}
	tests := []struct {
		name      string
		args      args
		want      entity.CurrentPrice
		wantErr   bool
	}{
		{"All good", args{response:goodResponse}, goodCurrentPrice, false},
		{"All good another date format", args{response:goodResponseAnotherDateFormat}, goodCurrentPrice, false},
		{"When wrong last price should return error", args{response:wrongLastPrice}, entity.CurrentPrice{}, true},
		{"When wrong last time should return error", args{response:wrongLastPriceTime}, entity.CurrentPrice{}, true},
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
