package getData

import (
	//"errors"
	"time"
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"errors"
)

type TimeSeries struct {
	StockName      string
	LastDate       string
	LastClosePrice float64
	DataName       string
	Dates          []time.Time
	ClosePrices    []float64
}

type QuandlClosePriceTimeSeries struct {
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
			PriceSeries         []QuandlPriceDayData `json:"data"`
		} `json:"dataset"`
}

type QuandlPriceDayData struct{
	Date string
	ClosingPrice float64
}

func (n *QuandlPriceDayData) UnmarshalJSON(buf []byte) error {
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

func PriceSeries(ticker string) (error, *TimeSeries) {
	url := generateUrl(ticker)
	rawTs := new(QuandlClosePriceTimeSeries)
	err := getJson(url, rawTs)
	if err == nil && rawTs.Dataset.Name == "" {
		err = errors.New("didn't get data for ticker")
		cleanTs := TimeSeries{}
		return err, &cleanTs
	}
	cleanTs := cleanUpTimeSeries(rawTs)
	return err, &cleanTs
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

func cleanUpTimeSeries(rawTs *QuandlClosePriceTimeSeries) (processedTs TimeSeries) {
	lastDay := rawTs.Dataset.PriceSeries[0]

	processedTs.StockName = rawTs.Dataset.Name
	processedTs.LastDate = lastDay.Date
	processedTs.LastClosePrice = lastDay.ClosingPrice
	processedTs.DataName = rawTs.Dataset.ColumnNames[1]
	processedTs.Dates, processedTs.ClosePrices = getXYvals(rawTs)
	return
}

func getXYvals(ts *QuandlClosePriceTimeSeries)  (xValues []time.Time, yValues []float64,) {
	for _, dayData := range ts.Dataset.PriceSeries {
		xDate, _ := time.Parse("2006-01-02", dayData.Date)
		xValues = append(xValues, xDate)
		yValues = append(yValues, dayData.ClosingPrice)
	}
	return
}

