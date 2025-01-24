package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	my_error "github.com/chyngyz-sydykov/go-web/error"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-web/internal/book"
)

type BookHandler struct {
	service       book.BookService
	commonHandler CommonHandler
}

type BookResponseDto struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	ICBN        string `json:"icbn"`
	PublishedAt time.Time
	AuthorId    uint `json:"author_id"`
}

func NewBookHandler(service book.BookService, commonHandler CommonHandler) *BookHandler {
	return &BookHandler{service: service, commonHandler: commonHandler}
}

// @Summary Get all books
// @Description Retrieve a list of all books
// @Tags books
// @Produce json
// @Success 200 {array} models.Book
// @Failure 500 {object} ErrorResponse
// @Router /books [get]
func (handler *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	books, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		handler.commonHandler.HandleError(w, err, http.StatusInternalServerError,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    SERVER_ERROR,
					Message: "server error",
				},
			})
		return
	}
	json.NewEncoder(w).Encode(books)
}

// @Summary Get a book by ID
// @Description Retrieve details of a specific book by its ID
// @Tags books
// @Produce json
// @Param bookId path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 400 {object} ErrorResponse "Invalid book ID"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /books/{bookId} [get]
func (handler *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	bookId, err := strconv.ParseUint(r.PathValue("bookId"), 10, 64)

	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    INVALID_REQUEST,
					Message: "invalid book id",
				},
			})
		return
	}

	book, err := handler.service.GetByID(int(bookId))
	if err != nil {
		if errors.Is(err, my_error.ErrgRpcServerDown) {
			handler.commonHandler.LogError(err, http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(book)
			return
		} else {
			handler.commonHandler.HandleError(w, err, http.StatusNotFound,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    RESOURCE_NOT_FOUND,
						Message: "Book with specified id is not found.",
					},
				})
		}
		return
	}
	json.NewEncoder(w).Encode(book)
}

// @Summary Create a new book
// @Description Add a new book to the library
// @Tags books
// @Accept json
// @Produce json
// @Param book body book.BookRequest true "Book payload"
// @Success 201 {object} models.Book
// @Failure 400 {object} ErrorResponse "Invalid payload"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /books [post]
func (handler *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.BookRequest
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    INVALID_REQUEST,
					Message: "provided book payload is invalid.",
				},
			})
		return
	}

	book, err := handler.service.Create(bookRequest)
	if err != nil {

		handler.commonHandler.HandleError(w, err, http.StatusInternalServerError,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    SERVER_ERROR,
					Message: "payload is corrupted.",
				},
			})
		return
	}
	w.Header().Set("Location", "api/v1/books/"+strconv.Itoa(int(book.ID)))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// @Summary Update a book
// @Description Update details of a specific book
// @Tags books
// @Accept json
// @Produce json
// @Param bookId path int true "Book ID"
// @Param book body book.BookRequest true "Book payload"
// @Success 200 {object} models.Book
// @Failure 400 {object} ErrorResponse "Invalid book ID or payload"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /books/{bookId} [put]
func (handler *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	bookId, err := strconv.ParseUint(r.PathValue("bookId"), 10, 64)

	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    INVALID_REQUEST,
					Message: "invalid book id",
				},
			})
		return
	}
	var bookRequest book.BookRequest

	err = json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    INVALID_REQUEST,
					Message: "provided payload is invalid.",
				},
			})
		return
	}
	var updatedBook *models.Book
	updatedBook, err = handler.service.Update(int(bookId), bookRequest)
	if err != nil {
		if errors.Is(err, my_error.ErrNotFound) {
			handler.commonHandler.HandleError(w, err, http.StatusNotFound,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    RESOURCE_NOT_FOUND,
						Message: "Book with specified id is not found.",
					},
				})
			return
		} else {
			handler.commonHandler.HandleError(w, err, http.StatusInternalServerError,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    SERVER_ERROR,
						Message: "Couldn't update payload.",
					},
				})
		}
		return
	}
	w.Header().Set("Location", "api/v1/books/"+strconv.Itoa(int(bookId)))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBook)
}

// Delete handles the HTTP request to delete a book by its ID.
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Produce json
// @Param bookId path int true "Book ID"
// @Success 200 {object} string "Book deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid book ID"
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/books/{bookId} [delete]

func (handler *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {

	bookId, err := strconv.ParseUint(r.PathValue("bookId"), 10, 64)

	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    INVALID_REQUEST,
					Message: "invalid book id",
				},
			})
		return
	}
	err = handler.service.Delete(int(bookId))

	if err != nil {
		if errors.Is(err, my_error.ErrNotFound) {
			handler.commonHandler.HandleError(w, err, http.StatusNotFound,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    RESOURCE_NOT_FOUND,
						Message: "Book with specified id is not found.",
					},
				})
			return
		} else {
			handler.commonHandler.HandleError(w, err, http.StatusInternalServerError,
				ErrorResponse{
					Error: ErrorDetail{
						Code:    SERVER_ERROR,
						Message: "Couldn't update payload.",
					},
				})
			return
		}

	}
	w.WriteHeader(http.StatusOK)
}
