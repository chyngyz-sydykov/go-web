package main

import (
	"log"
	"net/http"

	"github.com/chyngyz-sydykov/go-web/application"
	"github.com/chyngyz-sydykov/go-web/config"
	"github.com/chyngyz-sydykov/go-web/db"
	"github.com/chyngyz-sydykov/go-web/middleware"
	"github.com/chyngyz-sydykov/go-web/router"
	"gorm.io/gorm"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	db := initializeDatabase()
	app := application.InitializeApplication(db)
	router := router.InitializeRouter(app)

	middlewareController := middleware.NewMiddlewareController()

	http.ListenAndServe(
		":"+config.ApplicationPort,
		middlewareController.Chain()(router))
}

func initializeDatabase() *gorm.DB {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Could not load database config: %v", err)
	}
	dbInstance, err := db.InitializeDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	db.Migrate()
	return dbInstance

}
