package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
)

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
	var apiKey string = os.Getenv("QUANDL_API_KEY")
	resp, _ := http.Get("https://www.quandl.com/api/v3/datasets/WIKI/F.json?api_key=" + apiKey)
	defer resp.Body.Close()
	htmlData, _ := ioutil.ReadAll(resp.Body)

	fmt.Fprintf(w, string(htmlData))

}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/key", keyHandler)
	http.HandleFunc("/raw", rawHandler)
	http.ListenAndServe(":8080",nil)
}
