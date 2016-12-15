package priceChart

import (
	"github.com/wcharczuk/go-chart"
	"time"
	"bytes"
)

func GenerateChart(xSeries []time.Time, ySeries []float64) *bytes.Buffer {

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			Name:      "Date",
			NameStyle: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Name:      "Close Price/USD",
			NameStyle: chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xSeries,
				YValues: ySeries,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.SVG, buffer)
	return buffer
}
