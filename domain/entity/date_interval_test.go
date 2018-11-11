package entity

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
