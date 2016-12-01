package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
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
			RefreshedAt time.Time `json:"refreshed_at"`
			NewestAvailableDate string `json:"newest_available_date"`
			OldestAvailableDate string `json:"oldest_available_date"`
			ColumnNames []string `json:"column_names"`
			Frequency string `json:"frequency"`
			Type string `json:"type"`
			Premium bool `json:"premium"`
			Limit interface{} `json:"limit"`
			Transform interface{} `json:"transform"`
			ColumnIndex int `json:"column_index"`
			StartDate string `json:"start_date"`
			EndDate string `json:"end_date"`
			Data []struct {
				dateString string `json:"0"`
				closePrice float64 `json:"1"`
			} `json:"data"`
			Collapse interface{} `json:"collapse"`
			Order interface{} `json:"order"`
			DatabaseID int `json:"database_id"`
		} `json:"dataset"`
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func handler(w http.ResponseWriter, r *http.Request) {
	os.Setenv("FORD_PRICE", "14")
	var fordPrice string = os.Getenv("FORD_PRICE")
	fmt.Fprintf(w, "stock price of Ford motor is " + fordPrice)
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	fmt.Fprintf(w, "api key is " + apiKey)
}

func rawHandler(w http.ResponseWriter, r *http.Request) {
	var justDateAndClose = "column_index=4&"
	var apiFilter string = justDateAndClose + "start_date=2016-11-20&"
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	resp, _ := http.Get("https://www.quandl.com/api/v3/datasets/WIKI/F.json?" + apiFilter + "api_key=" + apiKey)
	defer resp.Body.Close()
	htmlData, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, string(htmlData))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	var justDateAndClose = "column_index=4&"
	var apiFilter string = justDateAndClose + "start_date=2016-11-20&"
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	var url string = "https://www.quandl.com/api/v3/datasets/WIKI/F.json?" + apiFilter + "api_key=" + apiKey
	fordSeries := new(TimeSeries)
	getJson(url, fordSeries)
	fmt.Fprintf(w, fordSeries.Dataset.Name)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/key", keyHandler)
	http.HandleFunc("/raw", rawHandler)
	http.HandleFunc("/json", jsonHandler)
	http.ListenAndServe(":8080",nil)
}
