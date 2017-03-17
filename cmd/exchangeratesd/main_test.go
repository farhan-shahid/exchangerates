package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestHandleEcb(t *testing.T) {
	var tests = []struct {
		params        url.Values
		ExpectedCode  int
		ExpectedResp  rateResp
		ExpectedError string
	}{
		{
			params:        url.Values{"from": {"USD"}, "to": {"EUR"}, "date": {"2017-03-02"}},
			ExpectedCode:  http.StatusOK,
			ExpectedResp:  rateResp{Rate: 0.95111},
			ExpectedError: "",
		},
		{
			params:        url.Values{"from": {"USD"}, "to": {"EUR"}, "date": {"20-03-02"}},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedResp:  rateResp{},
			ExpectedError: "incorrect date format, should be similar to 2016-03-28",
		},
		{
			params:        url.Values{"from": {}, "to": {"EUR"}, "date": {"2017-03-02"}},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedResp:  rateResp{},
			ExpectedError: `missing "from" URL parameter`,
		},
		{
			params:        url.Values{"from": {"USD"}, "to": {}, "date": {"2017-03-02"}},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedResp:  rateResp{},
			ExpectedError: `missing "to" URL parameter`,
		},
	}

	handler := http.HandlerFunc(handleEcb)

	for i, tt := range tests {
		req, err := http.NewRequest("GET", "/ecb?"+tt.params.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.ExpectedCode {
			t.Errorf("#%d failed: expected code=%v, got %v", i, tt.ExpectedCode, rr.Code)
		}
		if rr.Code != http.StatusOK {
			if strings.TrimSpace(rr.Body.String()) != tt.ExpectedError {
				t.Errorf(`#%d failed: expected error=%q,
				got %q`, i, tt.ExpectedError, rr.Body.String())
			}
			continue
		}

		var resp rateResp
		err = json.NewDecoder(rr.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(resp, tt.ExpectedResp) {
			t.Errorf("#%d failed: expected resp=%v, got %v", i, tt.ExpectedResp, resp)
		}
	}
}

func TestHandleGoogleFinance(t *testing.T) {
	var tests = []struct {
		params        url.Values
		ExpectedCode  int
		ExpectedResp  rateResp
		ExpectedError string
	}{
		{
			params:        url.Values{"from": {"USD"}, "to": {"EUR"}, "date": {"2017-03-02"}},
			ExpectedCode:  http.StatusOK,
			ExpectedResp:  rateResp{Rate: 0.95111},
			ExpectedError: "",
		},
		{
			params:        url.Values{"from": {"USD"}, "to": {"EUR"}, "date": {"20-03-02"}},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedResp:  rateResp{},
			ExpectedError: "incorrect date format, should be similar to 2016-03-28",
		},
		{
			params:        url.Values{"from": {}, "to": {"EUR"}, "date": {"2017-03-02"}},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedResp:  rateResp{},
			ExpectedError: `missing "from" URL parameter`,
		},
		{
			params:        url.Values{"from": {"USD"}, "to": {}, "date": {"2017-03-02"}},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedResp:  rateResp{},
			ExpectedError: `missing "to" URL parameter`,
		},
	}

	handler := http.HandlerFunc(handleGoogleFinance)

	for i, tt := range tests {
		req, err := http.NewRequest("GET", "/googlefinance?"+tt.params.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.ExpectedCode {
			t.Errorf("#%d failed: expected code=%v, got %v", i, tt.ExpectedCode, rr.Code)
		}
		if rr.Code != http.StatusOK {
			if strings.TrimSpace(rr.Body.String()) != tt.ExpectedError {
				t.Errorf(`#%d failed: expected error=%q,
				got %q`, i, tt.ExpectedError, rr.Body.String())
			}
			continue
		}

		var resp rateResp
		err = json.NewDecoder(rr.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)
		}

		// if !reflect.DeepEqual(resp, tt.ExpectedResp) {
		// 	t.Errorf("#%d failed: expected resp=%v, got %v", i, tt.ExpectedResp, resp)
		// }
	}
}
