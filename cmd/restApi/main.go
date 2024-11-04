package main

import (
	"log"
	"net/http"

	"github.com/chyngyz-sydykov/go-web/config"
	"github.com/chyngyz-sydykov/go-web/db"
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/middleware"
	"gorm.io/gorm"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	_ = initDb()

	mux := http.NewServeMux()
	middlewareController := middleware.NewMiddlewareController()
	HelloHandler := http.HandlerFunc(handlers.HelloHandler)
	mux.Handle("/", middlewareController.Chain(HelloHandler))
	http.ListenAndServe(":"+config.ApplicationPort, mux)
}

func initDb() *gorm.DB {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Could not load database config: %v", err)
	}
	dbInstance, err := db.InitDb(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	db.Migrate()
	return dbInstance

}
