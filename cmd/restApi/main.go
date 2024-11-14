package main

import (
	"log"
	"net/http"

	"github.com/chyngyz-sydykov/go-web/application"
	"github.com/chyngyz-sydykov/go-web/application/middleware"
	"github.com/chyngyz-sydykov/go-web/application/router"
	"github.com/chyngyz-sydykov/go-web/infrastructure/config"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	app := application.InitializeApplication()
	router := router.InitializeRouter(app)

	middlewareController := middleware.NewMiddlewareController()

	http.ListenAndServe(
		":"+config.ApplicationPort,
		middlewareController.Chain()(router))
}
