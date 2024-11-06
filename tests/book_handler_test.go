package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/chyngyz-sydykov/go-web/application"
	"github.com/chyngyz-sydykov/go-web/db/models"
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/chyngyz-sydykov/go-web/logger"
	"github.com/chyngyz-sydykov/go-web/router"
)

func (suite *IntegrationSuite) TestShouldReturnSuccessResponseAndAllBooksWhenCallingGetAllEndpoint() {
	// arrange
	publishedAt := time.Now()
	expectedBookModel := models.Book{Title: "John Doe", ICBN: "sdlfjskdflsdf234", PublishedAt: &publishedAt}
	suite.db.Omit("AuthorId").Create(&expectedBookModel)

	req := httptest.NewRequest("GET", "/api/v1/books", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

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

	suite.db.Unscoped().Delete(&models.Book{}, expectedBookModel.ID)
}

func (suite *IntegrationSuite) TestShouldReturnSuccessResponseAndReturnSingleBookWhenCallingGetByIdEndpoint() {
	// arrange
	publishedAt := time.Now()
	expectedBookModel := models.Book{Title: "John Doe", ICBN: "sdlfjskdflsdf234", PublishedAt: &publishedAt}
	suite.db.Omit("AuthorId").Create(&expectedBookModel)

	req := httptest.NewRequest("GET", "/api/v1/books/"+strconv.Itoa(int(expectedBookModel.ID)), nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)

	var resultBook models.Book
	err := json.NewDecoder(resp.Body).Decode(&resultBook)
	suite.NoError(err)

	suite.Suite.Assert().Equal("John Doe", resultBook.Title)
	suite.Suite.Assert().Equal("sdlfjskdflsdf234", resultBook.ICBN)

	suite.db.Unscoped().Delete(&models.Book{}, expectedBookModel.ID)
}

func (suite *IntegrationSuite) TestShouldReturnNotFoundResponseWithErrorMessageWhenCallingGetByIdEndpointWithNotExistingBookId() {
	// arrange
	req := httptest.NewRequest("GET", "/api/v1/books/999", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusNotFound, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal("RESOURCE_NOT_FOUND", errorResponse.Error.Code)
}
func (suite *IntegrationSuite) TestShouldReturnBadResponseResponseWithErrorMessageWhenCallingGetByIdEndpointWithInvalidBookId() {
	// arrange
	req := httptest.NewRequest("GET", "/api/v1/books/invalidId", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal("INVALID_REQUEST", errorResponse.Error.Code)
}

func (suite *IntegrationSuite) TestShouldReturnBadResponseResponseWithErrorMessageWhenCallingCreatingWithInvalidPayload() {
	// arrange
	req := httptest.NewRequest("POST", "/api/v1/books", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal("INVALID_REQUEST", errorResponse.Error.Code)
}

func provideDependencies(suite *IntegrationSuite) *application.App {
	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)
	bookService := book.NewBookService(suite.db)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &application.App{
		BookHandler: *bookHandler,
	}
	return app
}
