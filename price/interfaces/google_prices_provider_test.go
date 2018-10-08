package interfaces

import (
	"reflect"
	"testing"
	"time"

	"org.alex859/stockprices/domain"
)


func Test_filterPriceList(t *testing.T) {
	var may5, _ = time.Parse("02-01-2006", "05-05-2018")
	var may6, _ = time.Parse("02-01-2006", "06-05-2018")
	var june6, _ = time.Parse("02-01-2006", "06-06-2018")
	var july6, _ = time.Parse("02-01-2006", "06-07-2018")
	var july16, _ = time.Parse("02-01-2006", "16-07-2018")
	var august6, _ = time.Parse("02-01-2006", "06-08-2018")
	var august16, _ = time.Parse("02-01-2006", "16-08-2018")
	var may6pp = domain.PricePoint{Time:may6, Price:200}
	var june6pp = domain.PricePoint{Time:june6, Price:201}
	var july6pp = domain.PricePoint{Time:july6, Price:2001}
	var august6pp = domain.PricePoint{Time:august6, Price:2301}
	var priceList = domain.PriceList{may6pp, june6pp, july6pp, august6pp}
	type args struct {
		list domain.PriceList
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name string
		args args
		want domain.PriceList
	}{
		{"When empty list should return empty list", args{from:may5, to:june6, list:domain.PriceList{}}, domain.PriceList{}},
		{"When to before from return empty list", args{from:august6, to:may6, list:priceList}, domain.PriceList{}},
		{"When all in interval return full", args{from:may5, to:august16, list:priceList}, priceList},
		{"When all in interval left extreme included return full", args{from:may6, to:august16, list:priceList}, priceList},
		{"When all in interval right extreme included return full", args{from:may6, to:july16, list:priceList}, domain.PriceList{may6pp, june6pp, july6pp}},
		{"When some in interval return only in interval", args{from:may6, to:august6, list:priceList}, priceList},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterPriceList(tt.args.list, tt.args.from, tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterPriceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeToPeriod(t *testing.T) {
	var firstOctober2018, _ = time.Parse("02-01-2006", "01-10-2018")
	var firstDecember2018, _ = time.Parse("02-01-2006", "01-11-2018")
	var today = time.Now()
	var today10AM = time.Date(today.Year(), today.Month(), today.Day(), 10, 0, 0, 0, time.UTC)
	var today12PM = time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, time.UTC)
	var threeDaysAgo = today.AddDate(0, 0, -3)
	//var fiveDaysAgo = today.AddDate(0, 0, -5)
	var twentyDaysAgo = today.AddDate(0, 0, -20)
	var oneMonthAgo = today.AddDate(0, -1, 0)
	var twoMonthsAgo = today.AddDate(0, -2, 0)
	//var sixMonthsAgo = today.AddDate(0, -6, 0)
	var sevenMonthsAgo = today.AddDate(0, -7, 0)
	//var oneYearAgo = today.AddDate(-1, 0, 0)
	var twoYearsAgo = today.AddDate(-2, -1, 0)
	var fiveYearsAgo = today.AddDate(-5, 0, 0)
	var sevenYearsAgo = today.AddDate(-7, 0, 0)
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    Period
		wantErr bool
	}{
		{"When to before from should return error", args{from:firstDecember2018, to:firstOctober2018}, Unknown, true},
		{"When from and to in today should return OneDay", args{from: today10AM, to: today12PM}, OneDay, false},
		{"When from between one and five days ago should return FiveDay", args{from: threeDaysAgo, to: today}, FiveDays, false},
		// to be able to test this we should mock time.Now()
		//{"When from five days ago should return FiveDay", args{from: fiveDaysAgo, to: today}, FiveDays, false},
		{"When from between five days and one month should return OneMonth", args{from: twentyDaysAgo, to: threeDaysAgo}, OneMonth, false},
		//{"When from one month should return OneMonth", args{from: oneMonthAgo, to: threeDaysAgo}, OneMonth, false},
		{"When from between one month and six months ago should return SixMonth", args{from: twoMonthsAgo, to: oneMonthAgo}, SixMonth, false},
		//{"When from six months ago should return SixMonth", args{from: sixMonthsAgo, to: today}, SixMonth, false},
		{"When from between six months and one year ago should return OneYear", args{from: sevenMonthsAgo, to: oneMonthAgo}, OneYear, false},
		//{"When from one year ago should return OneYear", args{from: oneYearAgo, to: today}, OneYear, false},
		{"When from between one year and five years ago should return FiveYears", args{from: twoYearsAgo, to: oneMonthAgo}, FiveYears, false},
		//{"When from five years ago should return FiveYears", args{from: fiveYearsAgo, to: today}, FiveYears, false},
		{"When from more than five years ago should return Max", args{from: sevenYearsAgo, to: fiveYearsAgo}, Max, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timeToPeriod(tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeToPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got != tt.want {
					t.Errorf("timeToPeriod() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
