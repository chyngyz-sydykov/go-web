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
	"github.com/chyngyz-sydykov/go-web/application/handlers"
	"github.com/chyngyz-sydykov/go-web/application/router"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-web/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-web/internal/book"
)

func (suite *IntegrationSuite) TestGetAllEndpoint_ShouldReturnSuccessResponseAndAllBooks() {
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

func (suite *IntegrationSuite) TestGetByIdEndpoint_ShouldReturnSuccessResponseAndReturnSingleBook() {
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

func (suite *IntegrationSuite) TestGetByIdEndpoint_ShouldReturnNotFoundResponseWithErrorMessage_WithNotExistingBookId() {
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
func (suite *IntegrationSuite) TestGetByIdEndpoint_ShouldReturnBadResponseWithErrorMessage_WithInvalidBookId() {
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

func (suite *IntegrationSuite) TestCreateEndpoint_ShouldReturnBadResponseWithErrorMessage_WithEmptyPayload() {
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

func (suite *IntegrationSuite) TestCreateEndpoint_ShouldReturnBadResponseWithErrorMessage_WithInvalidPayload() {
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

func (suite *IntegrationSuite) TestCreateEndpoint_ShouldReturnCreatedResponse_WithValidPayload() {
	// arrange
	testAuthor := models.Author{Firstname: "John", Lastname: "Create Doe"}
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

	err = suite.db.Where("title = ? and ICBN = ?", "new test book", "test_icbn").First(&models.Book{}).Error
	suite.Suite.Assert().Nil(err)

	suite.db.Unscoped().Delete(&models.Book{}, actualBook.ID)
	suite.db.Unscoped().Delete(&models.Author{}, testAuthor.ID)
}

func (suite *IntegrationSuite) TestUpdateEndpoint_ShouldReturnBadResponseWithErrorMessage_WithInvalidBookId() {
	// arrange
	req := httptest.NewRequest("PUT", "/api/v1/books/invalidId", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal(handlers.INVALID_REQUEST, errorResponse.Error.Code)
}
func (suite *IntegrationSuite) TestUpdateEndpoint_ShouldReturnBadResponseWithErrorMessage_WithNotExistingBookIdAndEmptyPayload() {
	// arrange
	req := httptest.NewRequest("PUT", "/api/v1/books/999", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal(handlers.INVALID_REQUEST, errorResponse.Error.Code)
}

func (suite *IntegrationSuite) TestUpdateEndpoint_ShouldReturnNotFoundResponseWithErrorMessage_WithNotExistingBookIdAndValidPayload() {
	// arrange
	payload := fmt.Sprintf(`{
    "title": "update test book",
    "icbn": "test_icbn",
    "authorId": %d
	}`, 25)

	req := httptest.NewRequest("PUT", "/api/v1/books/999", strings.NewReader(payload))

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusNotFound, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal(handlers.RESOURCE_NOT_FOUND, errorResponse.Error.Code)
}

func (suite *IntegrationSuite) TestUpdateEndpoint_ShouldReturnOkResponseWithBody_WithValidBookIdAndValidPayload() {
	// arrange
	testAuthor := models.Author{Firstname: "John", Lastname: "Doe"}
	suite.db.Create(&testAuthor)

	testAuthor2 := models.Author{Firstname: "John", Lastname: "Simpson"}
	suite.db.Create(&testAuthor2)

	publishedAt := time.Now()
	testBook := models.Book{Title: "test book", ICBN: "123423ASDF", PublishedAt: &publishedAt, AuthorId: int64(testAuthor.ID)}
	suite.db.Create(&testBook)

	payload := fmt.Sprintf(`{
    "title": "update test book",
    "icbn": "updated_test_icbn",
    "authorId": %d
}`, testAuthor2.ID)

	req := httptest.NewRequest("PUT", "/api/v1/books/"+strconv.Itoa(int(testBook.ID)), strings.NewReader(payload))

	w := httptest.NewRecorder()

	app := provideDependencies(suite)
	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()

	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)

	var actualBook models.Book
	err := json.NewDecoder(resp.Body).Decode(&actualBook)
	suite.NoError(err)

	suite.Suite.Assert().Equal("update test book", actualBook.Title)
	suite.Suite.Assert().Equal("updated_test_icbn", actualBook.ICBN)
	suite.Suite.Assert().Equal(testAuthor2.ID, uint(actualBook.AuthorId))

	err = suite.db.Where("title = ? and ICBN = ?", "update test book", "updated_test_icbn").First(&models.Book{}).Error
	suite.Suite.Assert().Nil(err)

	suite.db.Unscoped().Delete(&models.Book{}, actualBook.ID)
	suite.db.Unscoped().Delete(&models.Author{}, testAuthor.ID)
	suite.db.Unscoped().Delete(&models.Author{}, testAuthor2.ID)
}
func (suite *IntegrationSuite) TestDeleteEndpoint_ShouldReturnBadResponseWithErrorMessage_WithInvalidBookId() {
	// arrange
	req := httptest.NewRequest("DELETE", "/api/v1/books/invalidId", nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&errorResponse)

	suite.Suite.Assert().Equal(handlers.INVALID_REQUEST, errorResponse.Error.Code)
}
func (suite *IntegrationSuite) TestDeleteEndpoint_ShouldReturnNotFoundResponseWithErrorMessage_WithNotExistingBookId() {
	// arrange
	req := httptest.NewRequest("DELETE", "/api/v1/books/999", nil)

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
func (suite *IntegrationSuite) TestDeleteEndpoint_ShouldReturnOkResponse_WithExistingBookId() {
	testAuthor := models.Author{Firstname: "John", Lastname: "Delete Doe"}
	suite.db.Create(&testAuthor)

	publishedAt := time.Now()
	testBook := models.Book{Title: "deleting book", ICBN: "123423ASDF", PublishedAt: &publishedAt, AuthorId: int64(testAuthor.ID)}
	suite.db.Create(&testBook)

	// arrange
	req := httptest.NewRequest("DELETE", "/api/v1/books/"+strconv.Itoa(int(testBook.ID)), nil)

	w := httptest.NewRecorder()

	app := provideDependencies(suite)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	bodyBytes, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Suite.Assert().Equal("", string(bodyBytes))

	err := suite.db.Where("title = ?", "deleting book").First(&models.Book{}).Error
	suite.Suite.Assert().NotNil(err)

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
