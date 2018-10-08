package domain

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func Test_EmptyHistory_THEN_ReturnError(t *testing.T) {
	date := dateAtMidnightUTC(2018, time.January, 20)
	history := PriceList{}
	result, err := history.On(date)
	assert.Equal(t, PricePoint{} , result)
	assert.NotNil(t, err)
}

func Test_OnePriceBeforeDate_THEN_ReturnThatPrice(t *testing.T) {
	date := dateAtMidnightUTC(2018, time.January, 22)
	history := PriceList{
		PricePoint{Time:time.Date(2018, time.January, 21, 10, 10, 10, 0, time.UTC)},
	}
	result, err := history.On(date)
	assert.Equal(t, history[0], result)
	assert.Nil(t, err)
}

func Test_OnePriceAfterDate_THEN_ReturnError(t *testing.T) {
	date := dateAtMidnightUTC(2018, time.January, 20)
	history := PriceList{
		PricePoint{Time:time.Date(2018, time.January, 28, 10, 10, 10, 0, time.UTC)},
	}
	result, err := history.On(date)
	assert.Equal(t, PricePoint{}, result)
	assert.NotNil(t, err)
}

func Test_MorePricesBeforeDate_THEN_ReturnLast(t *testing.T) {
	date := dateAtMidnightUTC(2018, time.January, 26)
	history := PriceList{
		PricePoint{Time:time.Date(2018, time.January, 19, 10, 10, 10, 0, time.UTC)},
		PricePoint{Time:time.Date(2018, time.January, 21, 10, 10, 10, 0, time.UTC)},
		PricePoint{Time:time.Date(2018, time.January, 25, 10, 10, 10, 0, time.UTC)},
	}
	result, err := history.On(date)
	assert.Equal(t, history[2], result)
	assert.Nil(t, err)
}

func Test_MorePricesAfterDate_THEN_ReturnError(t *testing.T) {
	date := dateAtMidnightUTC(2018, time.January, 18)
	history := PriceList{
		PricePoint{Time:time.Date(2018, time.January, 19, 10, 10, 10, 0, time.UTC)},
		PricePoint{Time:time.Date(2018, time.January, 21, 10, 10, 10, 0, time.UTC)},
		PricePoint{Time:time.Date(2018, time.January, 25, 10, 10, 10, 0, time.UTC)},
	}
	result, err := history.On(date)
	assert.Equal(t, PricePoint{}, result)
	assert.NotNil(t, err)
}

func Test_MorePricesSameDay_THEN_ReturnError(t *testing.T) {
	date := time.Date(2018, time.January, 21, 10, 15, 30, 0, time.UTC)
	history := PriceList{
		PricePoint{Time:time.Date(2018, time.January, 19, 10, 10, 10, 0, time.UTC)},
		PricePoint{Time:time.Date(2018, time.January, 21, 17, 10, 10, 0, time.UTC)},
		PricePoint{Time:time.Date(2018, time.January, 25, 10, 10, 10, 0, time.UTC)},
	}
	result, err := history.On(date)
	assert.Equal(t, history[1], result)
	assert.Nil(t, err)
}

func dateAtMidnightUTC(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

