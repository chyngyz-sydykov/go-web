package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloHandler)
	http.ListenAndServe(":8080", nil)
}
