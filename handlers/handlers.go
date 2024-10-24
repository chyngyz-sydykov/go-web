package handlers

import (
	"fmt"
	"net/http"
)

func HelloHandler(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprint(responseWriter, "hello world")
}
