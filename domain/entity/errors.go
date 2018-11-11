package entity

import "fmt"

// ErrNothingFound defines an error where no prices can be found for a ticker.
type ErrNothingFound struct {
	errStr string
}

func (err ErrNothingFound) Error() string {
	return err.errStr
}

// NewErrNothingFound creates a new ErrNothingFound error.
func NewErrNothingFound(ticker Ticker) ErrNothingFound{
	return ErrNothingFound{errStr:fmt.Sprintf("Unable to get prices for ticker:%s", ticker)}
}
