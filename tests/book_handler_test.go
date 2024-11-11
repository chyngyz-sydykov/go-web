package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/chyngyz-sydykov/go-web/application"
	"github.com/chyngyz-sydykov/go-web/db/models"
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/chyngyz-sydykov/go-web/logger"
	"github.com/chyngyz-sydykov/go-web/router"
)

func (suite *IntegrationSuite) TestShouldReturnSuccessResponseAndAllBooks_WhenCallingGetAllEndpoint() {
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

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	defer resp.Body.Close()

	suite.Suite.Contains(bodyString, "John Doe")
	suite.Suite.Contains(bodyString, "sdlfjskdflsdf234")

	suite.db.Unscoped().Delete(&models.Book{}, expectedBookModel.ID)
}

func (suite *IntegrationSuite) TestShouldReturnSuccessResponseAndReturnSingleBook_WhenCallingGetByIdEndpoint() {
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

func (suite *IntegrationSuite) TestShouldReturnNotFoundResponseWithErrorMessage_WhenCallingGetByIdEndpoint_WithNotExistingBookId() {
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
func (suite *IntegrationSuite) TestShouldReturnBadResponseWithErrorMessage_WhenCallingGetByIdEndpoint_WithInvalidBookId() {
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

func (suite *IntegrationSuite) TestShouldReturnBadResponseWithErrorMessage_WhenCreating_WithEmptyPayload() {
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

func (suite *IntegrationSuite) TestShouldReturnBadResponseWithErrorMessage_WhenCreating_WithInvalidPayload() {
	// arrange
	payload := `{"invalidField": "invalidValue"}`
	req := httptest.NewRequest("POST", "/api/v1/books", strings.NewReader(payload))

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal(handlers.SERVER_ERROR, errorResponse.Error.Code)

}

func (suite *IntegrationSuite) TestShouldReturnCreatedResponse_WhenCreating_WithValidPayload() {
	// arrange
	testAuthor := models.Author{Firstname: "John", Lastname: "Doe"}
	suite.db.Create(&testAuthor)

	payload := fmt.Sprintf(`{
    "title": "new test book",
    "icbn": "test_icbn",
    "authorId": %d
	}`, testAuthor.ID)

	req := httptest.NewRequest("POST", "/api/v1/books", strings.NewReader(payload))

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()

	// Assert
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var actualBook models.Book
	err := json.NewDecoder(resp.Body).Decode(&actualBook)
	suite.NoError(err)

	suite.Suite.Assert().Equal("new test book", actualBook.Title)
	suite.Suite.Assert().Equal("test_icbn", actualBook.ICBN)
	suite.Suite.Assert().Equal(testAuthor.ID, uint(actualBook.AuthorId))

	suite.db.Unscoped().Delete(&models.Book{}, actualBook.ID)
	suite.db.Unscoped().Delete(&models.Author{}, testAuthor.ID)
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
