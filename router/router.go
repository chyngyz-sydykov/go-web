package router

import (
	"net/http"

	"github.com/chyngyz-sydykov/go-web/application"
)

func InitializeRouter(app *application.App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", app.BookHandler.GetAll)
	mux.HandleFunc("GET /books/{bookId}", app.BookHandler.GetByID)
	mux.HandleFunc("POST /books", app.BookHandler.Create)
	mux.HandleFunc("PUT /books/{bookId}", app.BookHandler.Update)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))
	return v1
}
