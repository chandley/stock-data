package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"encoding/json"
	"github.com/stock-data/priceChart"
	"errors"
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

type successDataModel struct {
	name string
	lastDate string
	lastClosePrice float64
	yValsName string
	xVals []time.Time
	yVals []float64
}

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.ListenAndServe(":" + port,nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Path[1:]
	err, data := getData(ticker)

	if err != nil {
		showErrorPage(w, ticker)
		return
	}
	showSuccessPage(w, data)
}

func getData(ticker string) (error, *successDataModel) {
	var url string = generateUrl(ticker)
	ts := new(TimeSeries)
	err := getJson(url, ts)
	if err == nil && ts.Dataset.Name == "" {
		err = errors.New("didn't get data for ticker")
		data := successDataModel{}
		return err, &data
	}
	data := generateSuccessDataModel(ts)
	return err, &data
}

func showErrorPage(w http.ResponseWriter, ticker string) {
	fmt.Fprintf(w, "<h1>Could not get data for %v</h1>", ticker)
	fmt.Fprintf(w, "<h3>Please try again</h3>")
}

func generateSuccessDataModel(ts *TimeSeries) (data successDataModel) {
	lastDay := ts.Dataset.PriceSeries[0]

	data.name = ts.Dataset.Name
	data.lastDate = lastDay.Date
	data.lastClosePrice = lastDay.ClosingPrice
	data.yValsName = ts.Dataset.ColumnNames[1]
	data.xVals, data.yVals = getXYvals(ts)
	return
}

func showSuccessPage(w http.ResponseWriter, data *successDataModel) {
	fmt.Fprintf(w, "<h1>%v</h1>", data.name)
	fmt.Fprintf(w, "<p>Close price on %v was <b>%.2f</b></p>", data.lastDate, data.lastClosePrice)
	fmt.Fprintf(w, "<h3>%v price graph</h3>", data.yValsName)
	fmt.Fprintf(w, "<body>")
	fmt.Fprint(w, priceChart.GenerateChart(data.xVals, data.yVals))
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



