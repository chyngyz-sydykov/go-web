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
	authorId    uint `json:"author_id"`
}

func NewBookHandler(service book.BookService, commonHandler CommonHandler) *BookHandler {
	return &BookHandler{service: service, commonHandler: commonHandler}
}

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

func (handler *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
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
	if err := handler.service.Create(&book); err != nil {

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

	var payload models.Book

	err = json.NewDecoder(r.Body).Decode(&payload)
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
	updatedBook, err = handler.service.Update(int(bookId), payload)
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
