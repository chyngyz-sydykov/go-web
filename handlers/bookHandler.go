package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chyngyz-sydykov/go-web/db/models"
	"github.com/chyngyz-sydykov/go-web/internal/book"
)

type BookHandler struct {
	service       book.BookService
	commonHandler CommonHandler
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
					Code:    "SERVER_ERROR",
					Message: "server error",
				},
			})
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (handler *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	bookId, err := strconv.Atoi(r.PathValue("bookId"))

	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusBadRequest,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    "INVALID_REQUEST",
					Message: "invalid book id",
				},
			})
		return
	}

	book, err := handler.service.GetByID(bookId)
	if err != nil {
		handler.commonHandler.HandleError(w, err, http.StatusNotFound,
			ErrorResponse{
				Error: ErrorDetail{
					Code:    "RESOURCE_NOT_FOUND",
					Message: "Book with specified id is not found.",
				},
			})
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
					Code:    "INVALID_REQUEST",
					Message: "provided book payload is invalid.",
				},
			})
		return
	}
	if err := handler.service.Create(book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
