package main

import (
	"net/http"
	"fmt"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	os.Setenv("FORD_PRICE", "14")
	var fordPrice string = os.Getenv("FORD_PRICE")
	fmt.Fprintf(w, "stock price of Ford motor is " + fordPrice)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080",nil)
}
