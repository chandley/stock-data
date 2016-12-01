package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"encoding/json"
)

type TimeSeries struct {
	Dataset struct {
			ID int `json:"id"`
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
			Premium             bool `json:"premium"`
			Limit               interface{} `json:"limit"`
			Transform           interface{} `json:"transform"`
			ColumnIndex         int `json:"column_index"`
			StartDate           string `json:"start_date"`
			EndDate             string `json:"end_date"`
			PriceSeries         []PriceDayData `json:"data"`
			Collapse            interface{} `json:"collapse"`
			Order               interface{} `json:"order"`
			DatabaseID          int `json:"database_id"`
		} `json:"dataset"`
}

type PriceDayData struct{
	Date string
	ClosingPrice float64
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


func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}


func keyHandler(w http.ResponseWriter, r *http.Request) {
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	fmt.Fprintf(w, "api key is " + apiKey)
}

func generateUrl(ticker string) string {
	var justDateAndClose = "column_index=4&"
	var apiFilter string = justDateAndClose + "start_date=2016-11-20&"
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	return "https://www.quandl.com/api/v3/datasets/WIKI/" + ticker +".json?" + apiFilter + "api_key=" + apiKey
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	var url string = generateUrl("F")
	fordSeries := new(TimeSeries)
	getJson(url, fordSeries)
	fmt.Fprintf(w, fordSeries.Dataset.Name)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Close price was %f", fordSeries.Dataset.PriceSeries[0].ClosingPrice)
}

func main() {
	http.HandleFunc("/key", keyHandler)
	http.HandleFunc("/json", jsonHandler)
	http.ListenAndServe(":8080",nil)
}
