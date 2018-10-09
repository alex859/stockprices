package interfaces

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/Knetic/govaluate"
	"strconv"
	"time"
	"regexp"
	"reflect"
	"golang.org/x/net/html"
	"org.alex859/stockprices/domain"
)

type googlePricesClient struct {
	httpClient *http.Client
}

func NewGooglePrices(client *http.Client) *googlePricesClient {
	return &googlePricesClient{client}
}

func (gp *googlePricesClient) FetchPrices(ticker domain.Ticker, period Period) (*domain.PriceHistory, error) {
	symbol := fmt.Sprintf("%s:%s", ticker.Market, ticker.Symbol)
	defer timeTrack(time.Now(), symbol)
	symbolEncoded, eiCode, err := gp.searchSymbolAndEi(symbol)
	if err != nil {
		log.Printf("Unable to get symbol for ticker: %s. Error: %s", ticker, err)
		return nil, err
	}

	quotes, err := gp.getQuotesText(symbolEncoded, eiCode, period)

	if err != nil {
		log.Printf("Unable to get quotes for ticker: %s. Error: %s", ticker, err)
		return nil, err
	}

	return readGoogleResponse(quotes)
}

func (gp *googlePricesClient) getQuotesText(symbolEncoded, eiCode string, period Period) (string, error) {
	const quotesUrlTemplate = "https://www.google.com/async/finance_wholepage_chart?ei=%s&yv=3&async=mid_list:%s,period:%s,interval:%s,extended:true,element_id:fw-uid_%s_1,_id:fw-uid_%s_1,_pms:s,_fmt:pc"
	quotesText, err := gp.htmlFrom(fmt.Sprintf(quotesUrlTemplate, eiCode, symbolEncoded, period.Value().Str, period.Value().Interval, eiCode, eiCode))
	if err != nil {
		return "", err
	}

	return quotesText, nil
}

func (gp *googlePricesClient) searchSymbolAndEi(ticker string) (dataMid string, ei string, err error) {
	const searchTemplate = "https://www.google.com/search?hl=en&q=%s&btnG=Google+Search&tbs=0&safe=off&tbm=fin"
	resultsDocument, err := gp.document(gp.htmlFrom(fmt.Sprintf(searchTemplate, ticker)))
	if err != nil {
		return "", "", err
	}

	dataMid, found := gp.findFirstAttributeValue(resultsDocument.Find("#rso > div > div").Nodes, "data-mid")
	if !found {
		return "", "", fmt.Errorf("unable to find symbol: %s", ticker)
	}

	ei, found = gp.findFirstAttributeValue(resultsDocument.Find(`#tophf > input[name="ei"]`).Nodes, "value")
	if !found {
		return "", "", fmt.Errorf("unable to find ei for symbol: %s", ticker)
	}

	return dataMid, ei, nil
}

func (gp *googlePricesClient) findFirstAttributeValue(nodes []*html.Node, attrName string) (attrValue string, found bool) {
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == attrName {
				return attr.Val, true
			}
		}
	}

	return "", false
}

func (gp *googlePricesClient) document(html string, err error) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (gp *googlePricesClient) htmlFrom(url string) (string, error) {
	const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"User-Agent": []string{userAgent},
	}

	response, err := gp.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(htmlBytes), nil
}

func readGoogleResponse(str string) (*domain.PriceHistory, error) {
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
		return nil, fmt.Errorf("unexpected google response format reading: %s", str)
	}

	start := strings.Index(text, `"[`)
	end := strings.Index(text, `]\n"`)

	stringToAnalyse := text[start+1 : end+1]
	stringToAnalyse = strings.Replace(stringToAnalyse, `\n`, "", -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, `\"`, `"`, -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, "[", "(", -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, "]", ")", -1)
	stringToAnalyse = strings.Replace(stringToAnalyse, "null", `0`, -1)

	result, err := govaluate.NewEvaluableExpression(stringToAnalyse)
	if err != nil {
		return nil, err
	}
	r, err := result.Evaluate(nil)
	if err != nil {
		return nil, err
	}

	data, err := i2s(r)
	if err != nil {
		return nil, err
	}
	data2, err := i2s(data[2])
	if err != nil {
		return nil, err
	}
	p, err := strconv.ParseFloat(removeCommas(data2[2]), 64)
	if err != nil {
		return nil, err
	}
	prices := []domain.PricePoint{{Price: p, Time: time.Unix(toTimestamp(data2[4]), 0)}}
	for i := 5; i < len(data2)-5; i++ {
		d, err := i2s(data2[i])
		if err != nil {
			return nil, err
		}
		p, err := strconv.ParseFloat(removeCommas(d[2]), 64)
		if err != nil {
			return nil, err
		}
		prices = append(prices, domain.PricePoint{Price: p, Time: time.Unix(toTimestamp(d[4]), 0)})
	}

	data6, err := i2s(data[6])
	if err != nil {
		data6, err = i2s(data[7])
		if err != nil {
			return nil, err
		}
	}
	data617, err := i2s(data6[17])
	if err != nil {
		return nil, err
	}

	return &domain.PriceHistory{
		Currency:      data6[7].(string),
		Name:          data617[1].(string),
		Ticker:        domain.Ticker{Symbol:data617[2].(string), Market:data6[3].(string)},
		LastPrice:     data617[4].(float64),
		LastPriceTime: data617[8].(string),
		Prices:        prices,
	}, nil
}

func toTimestamp(data interface{}) int64 {
	return int64(60*data.(float64))
}

func removeCommas(str interface{}) string {
	return strings.Replace(str.(string), ",", "", -1)
}

func i2s(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("unexpected non slice type: %s", slice)
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}