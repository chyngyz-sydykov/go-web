package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/chyngyz-sydykov/go-web/application"
	"github.com/chyngyz-sydykov/go-web/db/models"
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/chyngyz-sydykov/go-web/router"
)

func (suite *IntegrationSuite) TestGetAllEndpoint() {
	// arrange
	publishedAt := time.Now()
	expectedBookModel := models.Book{Title: "John Doe", ICBN: "sdlfjskdflsdf234", PublishedAt: &publishedAt}
	suite.db.Omit("AuthorId").Create(&expectedBookModel)

	req := httptest.NewRequest("GET", "/api/v1/books", nil)

	w := httptest.NewRecorder()

	bookService := book.NewBookService(suite.db)

	bookHandler := handlers.NewBookHandler(*bookService)
	app := &application.App{
		BookHandler: *bookHandler,
	}
	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)

	var books []models.Book
	err := json.NewDecoder(resp.Body).Decode(&books)
	suite.NoError(err)

	resultBook := books[0]
	suite.Suite.Assert().Equal("John Doe", resultBook.Title)
	suite.Suite.Assert().Equal("sdlfjskdflsdf234", resultBook.ICBN)
}
