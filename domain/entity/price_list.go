package entity

import (
	"time"
	"sort"
	"github.com/pkg/errors"
)

type (
	// PricePoint defines a price at a given time.
	PricePoint struct {
		Price float64   `json:"price"`
		Time  time.Time `json:"time"`
	}

	// PriceList is a list of PricePoints.
	PriceList []PricePoint
)

// On gets the PricePoint at a given point in time.
func (pl PriceList) On(day time.Time) (PricePoint, error) {
	if len(pl) == 0 {
		return PricePoint{}, errors.New("unable to get price on date for empty price history")
	}

	sort.Slice(pl, func(i, j int) bool {
		return pl[i].Time.Before(pl[j].Time)
	})

	dayToSearch := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, time.UTC)
	var i int
	for i = 0; i < len(pl) && pl[i].Time.Before(dayToSearch); i++ {

	}

	if i == 0 {
		return PricePoint{}, errors.New("nothing found")
	}

	// we should check if the day param is too far ahead
	return pl[i-1], nil
}

// FilterByInterval returns a new PriceList containing only the PricePoints in the given DateInterval.
func (pl PriceList) FilterByInterval(interval DateInterval) PriceList {
	var result = PriceList{}
	for _, price := range pl {
		if interval.Contains(price.Time) {
			result = append(result, price)
		}
	}
	return result
}

