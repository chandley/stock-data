package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"encoding/json"
	"bytes"
	"github.com/wcharczuk/go-chart"
)

type TimeSeries struct {
	Dataset struct {
			DatasetCode string `json:"dataset_code"`
			DatabaseCode string `json:"database_code"`
			Name string `json:"name"`
			Description string `json:"description"`
			RefreshedAt         time.Time `json:"refreshed_at"`
			NewestAvailableDate string `json:"newest_available_date"`
			OldestAvailableDate string `json:"oldest_available_date"`
			ColumnNames         []string `json:"column_names"`
			Frequency           string `json:"frequency"`
			Type                string `json:"type"`
			ColumnIndex         int `json:"column_index"`
			StartDate           string `json:"start_date"`
			EndDate             string `json:"end_date"`
			PriceSeries         []PriceDayData `json:"data"`
		} `json:"dataset"`
}

type PriceDayData struct{
	Date string
	ClosingPrice float64
}

func main() {
	http.HandleFunc("/", jsonHandler)
	http.ListenAndServe(":8080",nil)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	var url string = generateUrl("FB")
	fetchedTimeSeries := new(TimeSeries)
	getJson(url, fetchedTimeSeries)

	fmt.Fprintf(w, "<h1>%v</h1>", fetchedTimeSeries.Dataset.Name)

	lastDay := fetchedTimeSeries.Dataset.PriceSeries[0]
	fmt.Fprintf(w, "<p>Close price on %v was <b>%.2f</b></p>", lastDay.Date, lastDay.ClosingPrice)

	fmt.Fprintf(w, "<h3>%v price graph</h3>", fetchedTimeSeries.Dataset.ColumnNames[1])

	fmt.Fprintf(w, "<body>")
	fmt.Fprint(w, generateChart(fetchedTimeSeries))
	fmt.Fprintf(w, "</body>")
}

func generateUrl(ticker string) string {
	var justDateAndClose = "column_index=4&"
	var apiFilter string = justDateAndClose
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	return "https://www.quandl.com/api/v3/datasets/WIKI/" + ticker +".json?" + apiFilter + "api_key=" + apiKey
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func (n *PriceDayData) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Date, &n.ClosingPrice}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Notification: %d != %d", g, e)
	}
	return nil
}

func getXYvals(ts *TimeSeries)  (xValues []time.Time, yValues []float64,) {
	for _, dayData := range ts.Dataset.PriceSeries {
		xDate, _ := time.Parse("2006-01-02", dayData.Date)
		xValues = append(xValues, xDate)
		yValues = append(yValues, dayData.ClosingPrice)
	}
	return
}

func generateChart(ts *TimeSeries) *bytes.Buffer {
	xSeries, ySeries := getXYvals(ts)
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
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

