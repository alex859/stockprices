package domain

type period struct {
	Str      string
	Interval string
}

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