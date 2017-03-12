package mockstore

import (
	"reflect"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	var tests = []struct {
		From         string
		To           string
		Date         string
		ExpectedErr  error
		ExpectedRate string
	}{
		{
			From:         "USD",
			To:           "EUR",
			Date:         "2017-03-02",
			ExpectedErr:  nil,
			ExpectedRate: "1",
		},
		{
			From:         "EUR",
			To:           "USD",
			Date:         "2017-03-02",
			ExpectedErr:  nil,
			ExpectedRate: "1",
		},
	}

	e := &MockStore{}
	for i, tt := range tests {
		got, err := e.GetExchangeRate(tt.From, tt.To, tt.Date)
		if want, got := tt.ExpectedErr, err; !reflect.DeepEqual(want, got) {
			t.Fatalf("#%d failed: expected error=%v, got %v", i, want, got)
		}
		if want, got := tt.ExpectedRate, got; want != got {
			t.Fatalf("#%d failed: expected rate=%v, got %v", i, want, got)
		}
	}
}
