package googlefinance

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"org.alex859/stockprices/domain"
	"reflect"
	"regexp"
	"strings"
)

// Service to go off to GoogleFinance and gather prices.
// Does not handle format errors, just returns a struct with string values to be validated later.
// TODO Google finance response parsing can be done in a decent way.
type defaultGoogleFinancePricesFetcher struct {
	httpClient *http.Client
}

// NewDefaultGooglePricesFetcher creates a new defaultGoogleFinancePricesFetcher.
func NewDefaultGooglePricesFetcher(client *http.Client) *defaultGoogleFinancePricesFetcher {
	return &defaultGoogleFinancePricesFetcher{client}
}

func (gp *defaultGoogleFinancePricesFetcher) FetchPrices(market string, symbol string, period domain.Period) (result domain.GoogleFinanceResponse, err error) {
	ticker := fmt.Sprintf("%s:%s", market, symbol)
	symbolEncoded, eiCode, err := gp.searchSymbolAndEi(ticker)
	if err != nil {
		message := fmt.Sprintf("unable to get symbol for ticker: %s. Error: %s", ticker, err)
		log.Print(message)
		err = errors.Wrap(err, message)
		return
	}

	quotes, err := gp.getQuotesText(symbolEncoded, eiCode, period)

	if err != nil {
		message := fmt.Sprintf("unable to get quotes for ticker: %s. Error: %s", ticker, err)
		log.Print(message)
		err = errors.Wrap(err, message)
		return
	}

	return readGoogleResponse(quotes)
}

func (gp *defaultGoogleFinancePricesFetcher) getQuotesText(symbolEncoded, eiCode string, period domain.Period) (string, error) {
	const quotesURLTemplate = "https://www.google.com/async/finance_wholepage_chart?ei=%s&yv=3&async=mid_list:%s,period:%s,interval:%s,extended:true,element_id:fw-uid_%s_1,_id:fw-uid_%s_1,_pms:s,_fmt:pc"
	quotesText, err := gp.htmlFrom(fmt.Sprintf(quotesURLTemplate, eiCode, symbolEncoded, period.Value().Str, period.Value().Interval, eiCode, eiCode))
	if err != nil {
		return "", errors.Wrap(err, "unable to read HTML")
	}

	return quotesText, nil
}

func (gp *defaultGoogleFinancePricesFetcher) searchSymbolAndEi(ticker string) (dataMid string, ei string, err error) {
	const searchTemplate = "https://www.google.com/search?hl=en&q=%s&btnG=Google+Search&tbs=0&safe=off&tbm=fin"
	resultsDocument, err := gp.document(gp.htmlFrom(fmt.Sprintf(searchTemplate, ticker)))
	if err != nil {
		return "", "", errors.Wrap(err, "unable to find symbol")
	}

	dataMid, found := gp.findFirstAttributeValue(resultsDocument.Find("#rso > div > div").Nodes, "data-mid")
	if !found {
		message := fmt.Sprintf("unable to find symbol: %s", ticker)
		return "", "", errors.New(message)
	}

	ei, found = gp.findFirstAttributeValue(resultsDocument.Find(`#tophf > input[name="ei"]`).Nodes, "value")
	if !found {
		message := fmt.Sprintf("unable to find ei for symbol: %s", ticker)
		return "", "", errors.New(message)
	}

	return dataMid, ei, nil
}

func (gp *defaultGoogleFinancePricesFetcher) findFirstAttributeValue(nodes []*html.Node, attrName string) (attrValue string, found bool) {
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == attrName {
				return attr.Val, true
			}
		}
	}

	return "", false
}

func (gp *defaultGoogleFinancePricesFetcher) document(html string, err error) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.Wrap(err, "erro querying html")
	}

	return document, nil
}

func (gp *defaultGoogleFinancePricesFetcher) htmlFrom(url string) (string, error) {
	const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create request")
	}
	req.Header = http.Header{
		"User-Agent": []string{userAgent},
	}

	response, err := gp.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to talk to remote server")
	}
	defer func() {
		if response.Body != nil {
			response.Body.Close()
		}
	}()

	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read response body")
	}

	return string(htmlBytes), nil
}

func readGoogleResponse(str string) (domain.GoogleFinanceResponse, error) {
	var financeResponse domain.GoogleFinanceResponse
	regex := regexp.MustCompile(`[0-9a-z]+;\[null`)
	text := ""
	found := false
	for _, row := range strings.Split(str, "\n") {
		if regex.MatchString(row) {
			found = true
		}
		if found {
			text = text + row
		}
	}

	if !found {
		message := fmt.Sprintf("unexpected google response format reading: %s", str)
		return financeResponse, errors.New(message)
	}

	start := strings.Index(text, `"[`)
	end := strings.Index(text, `]\n"`)

	stringToAnalyse := text[start+1 : end+1]
	stringToAnalyse = strings.Replace(stringToAnalyse, `\n`, "", -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, `\"`, `"`, -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, "[", "(", -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, "]", ")", -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, "null", `0`, -1)

	evaluableExpression, err := govaluate.NewEvaluableExpression(stringToAnalyse)
	if err != nil {
		return financeResponse, errors.Wrap(err, "unable to parse google finance response")
	}
	r, err := evaluableExpression.Evaluate(nil)
	if err != nil {
		return financeResponse, errors.Wrap(err, "unable to parse google finance response")
	}

	data, err := i2s(r)
	if err != nil {
		return financeResponse, errors.Wrap(err, "unable to parse google finance response")
	}
	data2, err := i2s(data[2])
	if err != nil {
		return financeResponse, errors.Wrap(err, "unable to parse google finance response")
	}
	prices := []domain.GoogleFinancePriceRow{{Price: toString(data2[2]), Time: toString(data2[4])}}
	for i := 5; i < len(data2)-5; i++ {
		d, err := i2s(data2[i])
		if err != nil {
			return financeResponse, errors.Wrap(err, "unable to parse google finance response")
		}

		prices = append(prices, domain.GoogleFinancePriceRow{Price: toString(d[2]), Time: toString(d[4])})
	}

	data6, err := i2s(data[6])
	if err != nil {
		data6, err = i2s(data[7])
		if err != nil {
			return financeResponse, errors.Wrap(err, "unable to parse google finance response")
		}
	}
	data617, err := i2s(data6[17])
	if err != nil {
		return financeResponse, errors.Wrap(err, "unable to parse google finance response")
	}

	return domain.GoogleFinanceResponse{
		Currency:      data6[7].(string),
		Name:          data617[1].(string),
		Market:        data6[3].(string),
		Symbol:        data617[2].(string),
		LastPrice:     toString(data617[4]),
		LastPriceTime: data617[8].(string),
		PricesRows:    prices,
	}, nil
}

func i2s(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		message := fmt.Sprintf("unexpected non slice type: %s", slice)
		return nil, errors.New(message)
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func toString(anything interface{}) string {
	return fmt.Sprintf("%v", anything)
}
