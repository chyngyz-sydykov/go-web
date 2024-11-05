package application

import (
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"gorm.io/gorm"
)

type App struct {
	BookHandler handlers.BookHandler
}

func InitializeApplication(db *gorm.DB) *App {
	bookService := book.NewBookService(db)
	bookHandler := handlers.NewBookHandler(*bookService)

	app := &App{
		BookHandler: *bookHandler,
	}
	return app
}
