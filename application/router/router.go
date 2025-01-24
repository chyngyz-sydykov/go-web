package router

import (
	"net/http"

	"github.com/chyngyz-sydykov/go-web/application"
	_ "github.com/chyngyz-sydykov/go-web/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Go Web API
// @version         1.0
// @description     This is a sample server for a bookstore API.
// @contact.name   You wish. no support whatsoever
// @contact.email  chyngys6@gmail.com
// @host      localhost:8000
// @BasePath  /api/v1
// @schemes   http
func InitializeRouter(app *application.App) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", app.BookHandler.GetAll)
	mux.HandleFunc("GET /books/{bookId}", app.BookHandler.GetByID)
	mux.HandleFunc("POST /books", app.BookHandler.Create)
	mux.HandleFunc("PUT /books/{bookId}", app.BookHandler.Update)
	mux.HandleFunc("DELETE /books/{bookId}", app.BookHandler.Delete)

	mux.HandleFunc("POST /ratings", app.RatingHandler.SaveRating)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	root := http.NewServeMux()
	root.Handle("/swagger/", httpSwagger.WrapHandler)
	root.Handle("/api/v1/", v1)

	return root
}
