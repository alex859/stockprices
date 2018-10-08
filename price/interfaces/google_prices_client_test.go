package interfaces

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"org.alex859/stockprices/domain"
)

func Test_readGoogleResponse_Year_AllGood(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_year_all_good")
	result, err := readGoogleResponse(string(str))

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, domain.Ticker{Market:"LON", Symbol:"ANP"}, result.Ticker)
	assert.Equal(t, "GBX", result.Currency)
	assert.Equal(t, 172, len(result.Prices))
}

func Test_readGoogleResponse_Day_AllGood(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_day_all_good")
	result, err := readGoogleResponse(string(str))

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, domain.Ticker{Market:"LON", Symbol:"ANP"}, result.Ticker)
	assert.Equal(t, "GBX", result.Currency)
	assert.Equal(t, 6, len(result.Prices))
}

func Test_readGoogleResponse_WrongPriceFormat(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_day_wrong_price_format")
	result, err := readGoogleResponse(string(str))

	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func Test_readGoogleResponse_WrongOrder(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_day_unexpected_order")
	result, err := readGoogleResponse(string(str))

	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func Test_readGoogleResponse_Unknown(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_unknown")
	result, err := readGoogleResponse(string(str))

	assert.NotNil(t, err)
	assert.Nil(t, result)
}