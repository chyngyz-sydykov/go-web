package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	pb "github.com/chyngyz-sydykov/go-web/proto/rating"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *IntegrationSuite) TestGetByIdEndpoint_ShouldReturnSuccessResponseAndLogErrorIfRatingGrpcServerIsDown() {
	// arrange
	publishedAt := time.Now()
	expectedBookModel := models.Book{Title: "John Doe", ICBN: "sdlfjskdflsdf234", PublishedAt: &publishedAt}
	suite.db.Omit("AuthorId").Create(&expectedBookModel)

	req := httptest.NewRequest("GET", "/api/v1/books/"+strconv.Itoa(int(expectedBookModel.ID)), nil)

	w := httptest.NewRecorder()

	app, mockLogger := provideDependenciesWithMockRatingServerBeingDown(suite, int(expectedBookModel.ID))

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()

	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)

	mockLogger.AssertCalled(suite.T(), "LogError", http.StatusServiceUnavailable, my_error.ErrgRpcServerDown)

	var resultBook models.Book

	err := json.NewDecoder(resp.Body).Decode(&resultBook)
	suite.NoError(err)

	suite.Suite.Assert().Equal("John Doe", resultBook.Title)
	suite.Suite.Assert().Equal("sdlfjskdflsdf234", resultBook.ICBN)

	suite.db.Unscoped().Delete(&models.Book{}, expectedBookModel.ID)
}

func (suite *IntegrationSuite) TestGetByIdEndpoint_ShouldReturnSuccessResponseWithRatingsIfRatingGrpcServerIsUp() {
	// arrange
	testAuthor := models.Author{Firstname: "TestName", Lastname: "TestLastName"}
	suite.db.Create(&testAuthor)

	publishedAt := time.Now()
	expectedBookModel := models.Book{
		Title:       "John Doe",
		ICBN:        "sdlfjskdflsdf234",
		PublishedAt: &publishedAt,
		AuthorId:    int64(testAuthor.ID),
	}
	suite.db.Create(&expectedBookModel)
	req := httptest.NewRequest("GET", "/api/v1/books/"+strconv.Itoa(int(expectedBookModel.ID)), nil)

	w := httptest.NewRecorder()

	app := provideDependenciesWithMockRatingServerBeingUp(suite, int(expectedBookModel.ID))

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()
	// Assert
	suite.Equal(http.StatusOK, resp.StatusCode)

	var resultBook book.BookDTO

	err := json.NewDecoder(resp.Body).Decode(&resultBook)
	suite.NoError(err)

	suite.Suite.Assert().Equal("John Doe", resultBook.Title)
	suite.Suite.Assert().Equal("sdlfjskdflsdf234", resultBook.ICBN)

	suite.Suite.Assert().Equal(int32(5), resultBook.Ratings[0].Rating)
	suite.Suite.Assert().Equal("comment for some hash", resultBook.Ratings[0].Comment)
	suite.Suite.Assert().Equal("some hash", resultBook.Ratings[0].RatingID)

	suite.Suite.Assert().Equal("TestName", resultBook.Author.Firstname)
	suite.Suite.Assert().Equal("TestLastName", resultBook.Author.Lastname)
	suite.Suite.Assert().Equal(testAuthor.ID, resultBook.Author.ID)
	suite.Suite.Assert().Equal(int32(testAuthor.ID), resultBook.AuthorID)

	suite.db.Unscoped().Delete(&models.Book{}, expectedBookModel.ID)
	suite.db.Unscoped().Delete(&models.Author{}, testAuthor.ID)
}

func (suite *IntegrationSuite) TestCreateRatingEndpoint_ShouldReturnBadResponse_WithInValidPayload() {
	// arrange
	payload := fmt.Sprintf(`{
    "bookId": %d,
    "comment": "comment for some hash",
    "rating": %d
	}`, 6, 6)

	req := httptest.NewRequest("POST", "/api/v1/ratings", strings.NewReader(payload))

	w := httptest.NewRecorder()

	app, mockLogger := provideDependenciesWithMockRatingServerReturnsInvalidError()

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()

	// Assert
	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	mockLogger.AssertCalled(suite.T(), "LogError", http.StatusBadRequest, my_error.ErrInvalidArgument)
}

func (suite *IntegrationSuite) TestCreateRatingEndpoint_ShouldReturnServerError_WhenGrpcIsRaisingInteralError() {
	// arrange
	payload := fmt.Sprintf(`{
    "bookId": %d,
    "comment": "comment for some hash",
    "rating": %d
	}`, 6, 6)

	req := httptest.NewRequest("POST", "/api/v1/ratings", strings.NewReader(payload))

	w := httptest.NewRecorder()

	app, mockLogger := provideDependenciesWithMockGrpcIsRaisingInternalError()

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()

	// Assert
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)

	mockLogger.AssertCalled(suite.T(), "LogError", http.StatusInternalServerError, my_error.ErrgRpcServerDown)
}
func (suite *IntegrationSuite) TestCreateRatingEndpoint_ShouldReturnCreatedResponse_WithValidPayload() {
	// arrange
	payload := fmt.Sprintf(`{
    "bookId": %d,
    "comment": "comment for some hash",
    "rating": %d
	}`, 1, 5)

	req := httptest.NewRequest("POST", "/api/v1/ratings", strings.NewReader(payload))

	w := httptest.NewRecorder()

	app := provideDependenciesWithMockRatingServerBeingUp(suite, 1)

	router := router.InitializeRouter(app)

	router.ServeHTTP(w, req)

	resp := w.Result()

	// Assert
	suite.Equal(http.StatusCreated, resp.StatusCode)

	var actioalRatingDTO rating.RatingDTO
	err := json.NewDecoder(resp.Body).Decode(&actioalRatingDTO)
	suite.NoError(err)

	suite.Suite.Assert().Equal(int32(1), actioalRatingDTO.BookID)
	suite.Suite.Assert().Equal("comment for some hash", actioalRatingDTO.Comment)
	suite.Suite.Assert().Equal(int32(5), actioalRatingDTO.Rating)
	suite.Suite.Assert().Equal("some hash", actioalRatingDTO.RatingID)
}

type RatingServiceMock struct {
	mock.Mock
}
type GrcpClientMock struct {
	mock.Mock
}

func (m *GrcpClientMock) SaveRating(ctx context.Context, in *pb.SaveRatingRequest, opts ...grpc.CallOption) (*pb.SaveRatingResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SaveRatingResponse), args.Error(1)
}

// RPC to retrieve all ratings for a specific book
func (m *GrcpClientMock) GetRatings(ctx context.Context, in *pb.GetRatingsRequest, opts ...grpc.CallOption) (*pb.GetRatingsResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetRatingsResponse), args.Error(1)
}
func (m *RatingServiceMock) GetByBookId(bookId int) ([]rating.RatingDTO, error) {
	args := m.Called(bookId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]rating.RatingDTO), args.Error(1)
}
func (m *RatingServiceMock) Create(ratingDTO *rating.RatingDTO) error {
	args := m.Called(ratingDTO)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func provideDependenciesWithMockRatingServerBeingDown(suite *IntegrationSuite, bookid int) (*application.App, *MockLogger) {
	mockLogger := new(MockLogger)

	// Set up the expectation for LogError
	mockLogger.On("LogError", 503, my_error.ErrgRpcServerDown).Return()

	commonHandler := handlers.NewCommonHandler(mockLogger)

	var ratingServiceMock RatingServiceMock
	ratingServiceMock.On("GetByBookId", bookid).Return(nil, my_error.ErrgRpcServerDown)

	bookService := book.NewBookService(suite.db, &ratingServiceMock)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &application.App{
		BookHandler: *bookHandler,
	}
	return app, mockLogger
}

func provideDependenciesWithMockRatingServerReturnsInvalidError() (*application.App, *MockLogger) {
	var grpcClientMock GrcpClientMock
	grpcClientMock.On("SaveRating", mock.Anything, mock.Anything).Return(nil, status.Error(codes.InvalidArgument, "resource value(s) is invalid"))
	mockLogger := new(MockLogger)
	mockLogger.On("LogError", http.StatusBadRequest, my_error.ErrInvalidArgument).Return()
	ratingService := rating.NewRatingService(&grpcClientMock, time.Duration(5)*time.Second)
	commonHandler := handlers.NewCommonHandler(mockLogger)
	ratingHandler := handlers.NewRatingHandler(ratingService, *commonHandler)

	app := &application.App{
		RatingHandler: *ratingHandler,
	}
	return app, mockLogger
}
func provideDependenciesWithMockGrpcIsRaisingInternalError() (*application.App, *MockLogger) {
	var grpcClientMock GrcpClientMock
	grpcClientMock.On("SaveRating", mock.Anything, mock.Anything).Return(nil, status.Error(codes.Internal, "internal error"))
	mockLogger := new(MockLogger)
	mockLogger.On("LogError", http.StatusInternalServerError, my_error.ErrgRpcServerDown).Return()
	ratingService := rating.NewRatingService(&grpcClientMock, time.Duration(5)*time.Second)
	commonHandler := handlers.NewCommonHandler(mockLogger)
	ratingHandler := handlers.NewRatingHandler(ratingService, *commonHandler)

	app := &application.App{
		RatingHandler: *ratingHandler,
	}
	return app, mockLogger
}

func provideDependenciesWithMockRatingServerBeingUp(suite *IntegrationSuite, bookId int) *application.App {
	logger := logger.NewLogger()
	commonHandler := handlers.NewCommonHandler(logger)

	var grpcClientMock GrcpClientMock

	getRatingRequest := &pb.GetRatingsRequest{
		BookId: int32(bookId),
	}
	getRatingsResponse := &pb.GetRatingsResponse{
		Ratings: []*pb.Rating{
			{
				RatingId: "some hash",
				BookId:   int32(bookId),
				Rating:   5,
				Comment:  "comment for some hash",
			}, {
				RatingId: "some hash 2",
				BookId:   int32(bookId),
				Rating:   3,
				Comment:  "comment for some hash 2",
			},
		}}

	saveRatingRequest := &pb.SaveRatingRequest{
		BookId:  int32(bookId),
		Rating:  5,
		Comment: "comment for some hash",
	}
	saveRatingResponse := &pb.SaveRatingResponse{
		Rating: &pb.Rating{
			RatingId: "some hash",
			BookId:   int32(bookId),
			Rating:   5,
			Comment:  "comment for some hash",
		},
	}

	grpcClientMock.On("GetRatings", mock.Anything, getRatingRequest).Return(getRatingsResponse, nil)

	grpcClientMock.On("SaveRating", mock.Anything, saveRatingRequest).Return(saveRatingResponse, nil)

	ratingService := rating.NewRatingService(&grpcClientMock, time.Duration(5)*time.Second)
	ratingHandler := handlers.NewRatingHandler(ratingService, *commonHandler)

	bookService := book.NewBookService(suite.db, ratingService)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &application.App{
		BookHandler:   *bookHandler,
		RatingHandler: *ratingHandler,
	}
	return app
}
