package main

import (
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Stock price of Ford motor is 10")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080",nil)
}
