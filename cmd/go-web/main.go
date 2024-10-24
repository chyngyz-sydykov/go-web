package main

import (
	"net/http"

	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/middlewares"
)

func main() {
	mux := http.NewServeMux()
	HelloHandler := http.HandlerFunc(handlers.HelloHandler)
	middlewares := middlewares.SetHeadersMiddleware(HelloHandler)
	mux.Handle("/", middlewares)
	http.ListenAndServe(":8080", mux)
}
