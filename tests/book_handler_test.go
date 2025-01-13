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
	my_error "github.com/chyngyz-sydykov/go-web/error"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-web/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/chyngyz-sydykov/go-web/internal/rating"
	"github.com/stretchr/testify/mock"
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

	app := provideDependenciesWithMessageBroker(suite, testBook.ID, strings.NewReader(payload))
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

	ratingDTO := &rating.RatingDTO{
		BookID:  6,
		Rating:  5,
		Comment: "comment for some hash",
	}

	var messageBrokerMock MessageBrokerMock
	messageBrokerMock.On("Publish", mock.Anything).Return(nil)

	var ratingServiceMock RatingServiceMock
	ratingServiceMock.On("GetByBookId", uint(6)).Return([]rating.RatingDTO{}, my_error.ErrgRpcServerDown)
	ratingServiceMock.On("Create", ratingDTO).Return(my_error.ErrgRpcServerDown)

	bookService := book.NewBookService(suite.db, &messageBrokerMock, &ratingServiceMock)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &application.App{
		BookHandler: *bookHandler,
	}
	return app
}

func provideDependenciesWithMessageBroker(suite *IntegrationSuite, bookId uint, body io.Reader) *application.App {
	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)

	ratingDTO := &rating.RatingDTO{
		BookID:  6,
		Rating:  5,
		Comment: "comment for some hash",
	}

	var messageBrokerMock MessageBrokerMock

	type Payload struct {
		Title    string `json:"title"`
		ICBN     string `json:"icbn"`
		AuthorId uint   `json:"authorId"`
	}

	// Parse the body to extract values
	var payload Payload
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse body: %v", err))
	}
	fmt.Println("Payload: ", payload)
	expectedMessage := book.BookMessage{
		ID:       bookId,
		Title:    payload.Title,
		ICBN:     payload.ICBN,
		EditedAt: time.Now(),
		Event:    "bookUpdated",
	}

	messageBrokerMock.On("Publish", mock.MatchedBy(func(msg book.BookMessage) bool {
		return msg.Title == expectedMessage.Title &&
			msg.ICBN == expectedMessage.ICBN &&
			msg.Event == expectedMessage.Event &&
			msg.ID == bookId
	})).Return(nil)

	var ratingServiceMock RatingServiceMock
	ratingServiceMock.On("GetByBookId", uint(6)).Return([]rating.RatingDTO{}, my_error.ErrgRpcServerDown)
	ratingServiceMock.On("Create", ratingDTO).Return(my_error.ErrgRpcServerDown)

	bookService := book.NewBookService(suite.db, &messageBrokerMock, &ratingServiceMock)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &application.App{
		BookHandler: *bookHandler,
	}
	return app
}

type BookPayload struct {
	Title    string `json:"title"`
	ICBN     string `json:"icbn"`
	AuthorID int64  `json:"authorId"`
}

// Function to read data from io.Reader and extract fields
func extractBookData(body io.Reader) (*BookPayload, error) {
	var payload BookPayload
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogError(statusCode int, err error) {
	m.Called(statusCode, err)
}
