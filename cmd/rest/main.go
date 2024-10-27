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

	router := http.NewServeMux()
	middlewareController := middleware.NewMiddlewareController()
	HelloHandler := http.HandlerFunc(handlers.HelloHandler)
	router.Handle("/", middlewareController.Chain(HelloHandler))

	http.ListenAndServe(":"+config.ApplicationPort, nil)
}
