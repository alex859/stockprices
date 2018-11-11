package googlefinance

import (
	"time"
	"org.alex859/stockprices/domain/entity"
)

type period struct {
	Str      string
	Interval string
}

// Period defines a period in the past we want to query GoogleFinance about.
type Period int
const (
	Max Period = iota
	FiveYears
	OneYear
	SixMonth
	OneMonth
	YearToDate
	FiveDays
	OneDay
)

// Value rte
func (p Period) Value() period {
	switch p {
	case Max:
		return period{Str: "40Y", Interval: "86400"}
	case FiveYears:
		return period{Str: "5Y", Interval: "86400"}
	case OneYear:
		return period{Str: "1Y", Interval: "86400"}
	case SixMonth:
		return period{Str: "6M", Interval: "86400"}
	case OneMonth:
		return period{Str: "1M", Interval: "86400"}
	case YearToDate:
		return period{Str: "YTD", Interval: "86400"}
	case FiveDays:
		return period{Str: "5d", Interval: "1800"}
	case OneDay:
		return period{Str: "1d", Interval: "300"}
	default:
		return period{}
	}
}

func (p Period) String() string {
	switch p {
	case Max:
		return "Max"
	case FiveYears:
		return "Five Years"
	case OneYear:
		return "One Year"
	case SixMonth:
		return "Six Months"
	case OneMonth:
		return "One Month"
	case YearToDate:
		return "Year To Date"
	case FiveDays:
		return "Five Days"
	case OneDay:
		return "One Day"
	default:
		return "Unknown"
	}
}

// FromDateInterval works out the period needed to get prices in the given date interval.
// E.g: If today is 1/12, and the time interval is 1/11 to 20/11, the required period will be 1 month.
func FromDateInterval(interval entity.DateInterval) Period {

	today := time.Now()
	sameDay := interval.From().Year() == today.Year() && interval.From().Month() == today.Month() && interval.From().Day() == today.Day()
	diff := time.Since(interval.From())

	day := time.Hour * 24
	month := day * 30
	year := day * 365
	switch {
	case sameDay:
		return OneDay
	case diff >= 5*year:
		return Max
	case diff >= 1*year:
		return FiveYears
	case diff >= 6*month:
		return OneYear
	case diff >= 1*month:
		return SixMonth
	case diff >= 5*day:
		return OneMonth
	default:
		return FiveDays
	}
}