package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chyngyz-sydykov/go-web/db/models"
	"github.com/chyngyz-sydykov/go-web/internal/book"
)

type BookHandler struct {
	service book.BookService
}

func NewBookHandler(service book.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (handler *BookHandler) swagger() {
}

func (handler *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	books, err := handler.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (handler *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	stringId := r.PathValue("bookId")
	if stringId == "" {
		http.Error(w, "book id not provided", http.StatusBadRequest)
		return
	}
	bookId, err := strconv.Atoi(stringId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	book, err := handler.service.GetByID(bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (handler *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := handler.service.Create(book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
