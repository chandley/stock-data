package main

import (
	"net/http"
	"fmt"
	"os"
	"github.com/stock-data/priceChart"
	"github.com/stock-data/getData"
	"github.com/stock-data/getFundamentals"
)

func main() {
	http.HandleFunc("/stock/", handler)
	http.HandleFunc("/fundamentals/", fundamentalsHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.ListenAndServe(":" + port,nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Path[len("/stock/"):]
	err, closePriceTimeSeries := getData.PriceSeries(ticker)

	if err != nil {
		showErrorPage(w, ticker, err)
		return
	}
	showPricePage(w, closePriceTimeSeries)
}

func fundamentalsHandler(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Path[len("/fundamentals/"):]
	fmt.Println(ticker)
	_, data := getFundamentals.LatestFundamentals(ticker)
	fmt.Fprintf(w, "<p>Company %v</p>", data["companyname"])
	fmt.Fprintf(w, "<p>Total assets %v</p>", data["totalassets"])
	fmt.Fprintf(w, "<p>Total revenue %v</p>", data["totalrevenue"])
	fmt.Fprintf(w, "<p>Net income %v</p>", data["netincome"])
}

func showErrorPage(w http.ResponseWriter, ticker string, err error) {
	fmt.Fprintf(w, "<h1>Could not get data for %v</h1>", ticker)
	fmt.Fprintf(w, "<h3>Please try again</h3>")
	fmt.Fprintf(w, "<p>Detected error %v</p>", err)

}

func showPricePage(w http.ResponseWriter, ts *getData.TimeSeries) {
	fmt.Fprintf(w, "<h1>%v</h1>", ts.StockName)
	fmt.Fprintf(w, "<p>Close price on %v was <b>%.2f</b></p>", ts.LastDate, ts.LastClosePrice)
	fmt.Fprintf(w, "<h3>%v price graph</h3>", ts.DataName)
	fmt.Fprintf(w, "<body>")
	fmt.Fprint(w, priceChart.GenerateChart(ts.Dates, ts.ClosePrices))
	fmt.Fprintf(w, "</body>")
}





