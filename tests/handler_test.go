package main

import (
	"log"
	"testing"

	"github.com/chyngyz-sydykov/go-web/infrastructure/config"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type IntegrationSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *IntegrationSuite) SetupSuite() {
	db := initializeDatabase()
	suite.db = db.Debug()
}

// TearDownSuite is called once after the test suite runs.
func (suite *IntegrationSuite) TearDownSuite() {
	// Clean up
	// err := suite.db.Exec("Truncate TABLE books;").Error
	// suite.Require().NoError(err, "Error dropping test table")

}

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func initializeDatabase() *gorm.DB {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Could not load database config: %v\n", err)
	}

	// Connect to the target database using Gorm
	dbInstance, err := db.InitializeDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	db.Migrate()
	return dbInstance
}
