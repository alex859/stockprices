package domain

import (
	"time"
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
)

func Test_FilterByInterval(t *testing.T) {
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var may6, _ = time.Parse("02-01-2006", "06-05-2018")
	var june6, _ = time.Parse("02-01-2006", "06-06-2018")
	var july6, _ = time.Parse("02-01-2006", "06-07-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	var august6, _ = time.Parse("02-01-2006", "06-08-2018")
	var august16, _ = time.Parse("02-01-2006", "16-08-2018")
	var may6pp = PricePoint{Time:may6, Price:200}
	var june6pp = PricePoint{Time:june6, Price:201}
	var july6pp = PricePoint{Time:july6, Price:2001}
	var august6pp = PricePoint{Time:august6, Price:2301}
	var priceList = PriceList{may6pp, june6pp, july6pp, august6pp}
	type args struct {
		list PriceList
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name string
		args args
		want PriceList
	}{
		{"When empty list should return empty list", args{from:may5, to:june6, list:PriceList{}}, PriceList{}},
		{"When all in interval return full", args{from:may5, to:august16, list:priceList}, priceList},
		{"When all in interval left extreme included return full", args{from:may6, to:august16, list:priceList}, priceList},
		{"When all in interval right extreme included return full", args{from:may6, to:july16, list:priceList}, PriceList{may6pp, june6pp, july6pp}},
		{"When some in interval return only in interval", args{from:may6, to:august6, list:priceList}, priceList},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interval, err := NewDateInterval(tt.args.from, tt.args.to)
			if assert.NoError(t, err) {
				if got := tt.args.list.FilterByInterval(interval); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("filterPriceList() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
