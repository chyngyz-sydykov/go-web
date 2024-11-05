package handlers

import (
	"fmt"
	"net/http"
)

type AppHandlerInterface interface {
	swagger()
}

func HelloHandler(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprint(responseWriter, "hello world")
}
