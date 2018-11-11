package handlers

import (
	"org.alex859/stockprices/domain/entity"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var noRequestParams = map[string]string{}

var fromInvalid = map[string]string{"from": "2014:15:62"}
var fromValid = map[string]string{"from": "12-12-2005"}
var dec12, _ = time.Parse("02-01-2006", "12-12-2005")
func Test_fromDate(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"No request params will return error", args{request(noRequestParams)}, time.Now(), true},
		{"Invalid params will return error", args{request(fromInvalid)}, time.Now(), true},
		{"Valid params will return correct time", args{request(fromValid)}, dec12 , false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromDate(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("fromDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

var toInvalid = map[string]string{"to": "3216aa"}
var toValid = map[string]string{"to": "15-10-2018"}
var oct15, _ = time.Parse("2-1-2006", "15-10-2018")
func Test_toDate(t *testing.T) {
	now = func() time.Time {
		return oct15
	}
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"No request params will return error", args{request(noRequestParams)}, oct15, false},
		{"Invalid params will return error", args{request(toInvalid)}, time.Now(), true},
		{"Valid params will return correct time", args{request(toValid)}, oct15 , false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToDate(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

var tickersInvalid = map[string]string{"tickers": "LON,SQ,NYSE"}
var tickersValid = map[string]string{"tickers": "LON:ANP,NYSE:SQ"}
var oneTickerValid = map[string]string{"tickers": "LON:ANP"}
var tickersSomeValid = map[string]string{"tickers": "LON:ANP,NYSESQ"}
var tickersValidWithSpaces1 = map[string]string{"tickers": "LON:ANP,NYSE: SQ"}
var tickersValidWithSpaces2 = map[string]string{"tickers": "LON:ANP, NYSE: SQ"}
var tickersNoComma = map[string]string{"tickers": "LON:ANP NYSE: SQ"}
func Test_tickers(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []entity.Ticker
		wantErr bool
	}{
		{"No request params will return error", args{request(noRequestParams)}, []entity.Ticker{}, true},
		{"Invalid tickers will return error", args{request(tickersInvalid)}, []entity.Ticker{}, true},
		{"Tickers will no comma return error", args{request(tickersNoComma)}, []entity.Ticker{}, true},
		{"Some valid tickers will return correct ones", args{request(tickersSomeValid)}, []entity.Ticker{{Symbol:"ANP", Market:"LON"}}, false},
		{"One valid ticker will return result", args{request(oneTickerValid)}, []entity.Ticker{{Symbol:"ANP", Market:"LON"}}, false},
		{"All valid tickers will return result", args{request(tickersValid)}, []entity.Ticker{{"LON", "ANP"}, {"NYSE", "SQ"}}, false},
		{"Tickers with spaces in ticker will return correct result", args{request(tickersValidWithSpaces1)}, []entity.Ticker{{"LON", "ANP"}, {"NYSE", "SQ"}}, false},
		{"Tickers with spaces between will return correct result", args{request(tickersValidWithSpaces2)}, []entity.Ticker{{"LON", "ANP"}, {"NYSE", "SQ"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tickers(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("tickers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tickers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func request(params map[string]string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{QueryStringParameters:params}
}