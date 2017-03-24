package chart

import (
	"bytes"
	"image"
	"image/gif"
	"io"
	"time"

	"github.com/farhan-shahid/exchangerates"
	gochart "github.com/wcharczuk/go-chart"
)

// MakeRateChartGIF writes a GIF graphical representation of the exchange rate data to the io.Writer provided
func MakeRateChartGIF(from, to string, rates []exchangerates.DateRate, w io.Writer) error {
	finalGIF := &gif.GIF{}
	var b bytes.Buffer

	for i := 2; i < len(rates); i++ {
		err := makeRateChartPartial(from, to, rates, &b, i+1)
		if err != nil {
			return err
		}

		pngImg, _, err := image.Decode(&b)
		if err != nil {
			return err
		}
		b.Reset()

		err = gif.Encode(&b, pngImg, &gif.Options{NumColors: 64})
		if err != nil {
			return err
		}

		gifImg, err := gif.Decode(&b)
		if err != nil {
			return err
		}

		finalGIF.Image = append(finalGIF.Image, gifImg.(*image.Paletted))
		finalGIF.Delay = append(finalGIF.Delay, 65)

		b.Reset()
	}
	finalGIF.Delay[len(finalGIF.Delay)-1] = 200
	return gif.EncodeAll(w, finalGIF)
}

// MakeRateChart writes a PNG graphical representation of the exchange rate data to the io.Writer provided
func MakeRateChart(from, to string, rates []exchangerates.DateRate, w io.Writer) error {
	return makeRateChartPartial(from, to, rates, w, len(rates))
}

func makeRateChartPartial(from, to string, rates []exchangerates.DateRate, w io.Writer, length int) error {
	dates := make([]time.Time, length)
	vals := make([]float64, length)

	for i := 0; i < length; i++ {
		dates[i] = rates[i].Date
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
					StrokeColor: gochart.ColorBlack.WithAlpha(128),
					//FillColor:   gochart.GetDefaultColor(0).WithAlpha(128),
				},
				XValues: dates,
				YValues: vals,
			},
		},
	}
	return graph.Render(gochart.PNG, w)
}
