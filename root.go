package exchangerates

import (
	"io"
	"time"

	gochart "github.com/wcharczuk/go-chart"
)

// Store is the interface that all exchange rate stores must satisfy
type Store interface {
	GetExchangeRate(from, to string, date string) (float64, error)
	GetMonthExchangeRates(from, to string, year, month int) ([]DateExchangeRatePair, error)
}

// DateExchangeRatePair represents a single exchange rate with its date
type DateExchangeRatePair struct {
	Date time.Time
	Rate float64
}

// MakeRateChart writes a graphical representation of the exchange rate data to io.Writer provided
func MakeRateChart(from, to string, rates []DateExchangeRatePair, w io.Writer) error {
	dates := make([]time.Time, len(rates))
	vals := make([]float64, len(rates))
	for i := range rates {
		dates[i] = rates[i].Date
	}
	for i := range rates {
		vals[i] = rates[i].Rate
	}

	graph := gochart.Chart{
		XAxis: gochart.XAxis{
			Style: gochart.StyleShow(),
		},
		YAxis: gochart.YAxis{
			Name:      from + " to " + to + " exchange rate",
			NameStyle: gochart.StyleShow(),
			Style:     gochart.StyleShow(),
		},
		Series: []gochart.Series{
			gochart.TimeSeries{
				Style: gochart.Style{
					Show:        true,
					StrokeColor: gochart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   gochart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: dates,
				YValues: vals,
			},
		},
	}
	return graph.Render(gochart.PNG, w)
}
