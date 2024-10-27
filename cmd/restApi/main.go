package main

import (
	"log"
	"net/http"

	"github.com/chyngyz-sydykov/go-web/config"
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/middleware"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	mux := http.NewServeMux()
	middlewareController := middleware.NewMiddlewareController()
	HelloHandler := http.HandlerFunc(handlers.HelloHandler)
	mux.Handle("/", middlewareController.Chain(HelloHandler))
	http.ListenAndServe(":"+config.ApplicationPort, mux)
}
