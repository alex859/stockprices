package domain

import (
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_NewInterval(t *testing.T) {
	var firstOctober2018, _ = time.Parse("02-01-2006", "01-10-2018")
	var firstDecember2018, _ = time.Parse("02-01-2006", "01-11-2018")
	_, err := NewDateInterval(firstDecember2018, firstOctober2018)
	assert.Error(t, err)

	interval, err := NewDateInterval(firstOctober2018, firstDecember2018)
	if assert.NoError(t, err) {
		assert.Equal(t, DateInterval{from: firstOctober2018, to: firstDecember2018}, interval)
	}
}

func Test_ToPeriod(t *testing.T) {
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
			interval, err := NewDateInterval(tt.args.from, tt.args.to)
			if assert.NoError(t, err) {
				got := interval.ToPeriod()
				if got != tt.want {
					t.Errorf("timeToPeriod() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
