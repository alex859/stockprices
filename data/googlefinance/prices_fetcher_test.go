package googlefinance

import (
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_readGoogleResponse_Year_AllGood(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_year_all_good")
	result, err := readGoogleResponse(string(str))

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, "LON", result.Market)
	assert.Equal(t, "ANP", result.Symbol)
	assert.Equal(t, "GBX", result.Currency)
	assert.Equal(t, 172, len(result.PricesRows))
}

func Test_readGoogleResponse_Day_AllGood(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_day_all_good")
	result, err := readGoogleResponse(string(str))

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, "LON", result.Market)
	assert.Equal(t, "ANP", result.Symbol)
	assert.Equal(t, "GBX", result.Currency)
	assert.Equal(t, 6, len(result.PricesRows))
}

func Test_readGoogleResponse_WrongOrder(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_day_unexpected_order")
	_, err := readGoogleResponse(string(str))

	assert.Error(t, err)
}

func Test_readGoogleResponse_Unknown(t *testing.T) {
	str, _ := ioutil.ReadFile("testdata/google_finance_unknown")
	_, err := readGoogleResponse(string(str))

	assert.Error(t, err)
}

