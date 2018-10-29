package domain

import (
	"time"
	"fmt"
	"github.com/pkg/errors"
)

type DateInterval struct {
	from time.Time
	to time.Time
}

func NewDateInterval(from time.Time, to time.Time) (DateInterval, error) {
	if to.Before(from) {
		return DateInterval{}, errors.New("parameter to must be after parameter from")
	}

	return DateInterval{from:from, to:to}, nil
}

// Works out the period of time containing the time interval.
// E.g: If today is 1/12, and the time interval is 1/11 to 20/11, the required period will be 1 month.
func (interval DateInterval) ToPeriod() Period {

	today := time.Now()
	sameDay := interval.from.Year() == today.Year() && interval.from.Month() == today.Month() && interval.from.Day() == today.Day()
	diff := time.Since(interval.from)

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

func (interval DateInterval) Contains(t time.Time) bool {
	return (t.After(interval.from) || t.Equal(interval.from)) && (t.Before(interval.to) || t.Equal(interval.to))
}

func (interval DateInterval) String() string {
	return fmt.Sprintf("%v to %v", interval.from, interval.to)
}