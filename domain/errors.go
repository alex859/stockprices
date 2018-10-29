package domain

import "fmt"

type ErrNothingFound struct {
	errStr string
}
func (err ErrNothingFound) Error() string {
	return err.errStr
}
func NewErrNothingFound(ticker Ticker) ErrNothingFound{
	return ErrNothingFound{errStr:fmt.Sprintf("Unable to get prices for ticker:%s", ticker)}
}
