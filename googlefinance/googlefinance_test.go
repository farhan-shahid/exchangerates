package googlefinance

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	var tests = []struct {
		From         string
		To           string
		Date         string
		ExpectedErr  error
		ExpectedRate float64
	}{
		{
			From:         "USD",
			To:           "EUR",
			Date:         "2017-03-02",
			ExpectedErr:  nil,
			ExpectedRate: 0.9399,
		},
		{
			From:         "EUR",
			To:           "USD",
			Date:         "2017-03-02",
			ExpectedErr:  nil,
			ExpectedRate: 1.0641,
		},
		{
			From:         "INR",
			To:           "USD",
			Date:         "2017-03-02",
			ExpectedErr:  nil,
			ExpectedRate: 0.01498,
		},
		{
			From:         "USD",
			To:           "XYZ",
			Date:         "2017-03-02",
			ExpectedErr:  errors.New("currency USD or XYZ not found"),
			ExpectedRate: 0,
		},
		{
			From:         "XYZ",
			To:           "USD",
			Date:         "2017-03-02",
			ExpectedErr:  errors.New("currency XYZ or USD not found"),
			ExpectedRate: 0,
		},
	}

	s := New()
	for i, tt := range tests {
		got, err := s.GetExchangeRate(tt.From, tt.To, tt.Date)
		if want, got := tt.ExpectedErr, err; !reflect.DeepEqual(want, got) {
			t.Fatalf("#%d failed: expected error=%v, got %v", i, want, got)
		}
		if want, got := tt.ExpectedRate, got; int(want) != int(got) { // matching integer parts only because rates change frequently
			t.Fatalf("#%d failed: expected rate=%v, got %v", i, want, got)
		}
	}
}
