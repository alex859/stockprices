package entity

import (
	"time"
	"fmt"
	"github.com/pkg/errors"
)

// DateInterval defines a date interval with a from and to date.
type DateInterval struct {
	from time.Time
	to time.Time
}

func (di DateInterval) From() time.Time {
	return di.from
}

func (di DateInterval) To() time.Time {
	return di.to
}

// NewDateInterval creates a new DateInterval for the given from and to times.
func NewDateInterval(from time.Time, to time.Time) (DateInterval, error) {
	if to.Before(from) {
		return DateInterval{}, errors.New("parameter to must be after parameter from")
	}

	return DateInterval{from:from, to:to}, nil
}

// Contains checks if a DateInterval contains the provided time.
func (interval DateInterval) Contains(t time.Time) bool {
	return (t.After(interval.from) || t.Equal(interval.from) || sameDay(t, interval.from)) && (t.Before(interval.to) || t.Equal(interval.to) || sameDay(t, interval.to))
}

func sameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year()
}

func (interval DateInterval) String() string {
	return fmt.Sprintf("%v to %v", interval.from, interval.to)
}