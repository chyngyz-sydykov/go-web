package application

import (
	"github.com/chyngyz-sydykov/go-web/handlers"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/chyngyz-sydykov/go-web/logger"
	"gorm.io/gorm"
)

type App struct {
	BookHandler handlers.BookHandler
}

func InitializeApplication(db *gorm.DB) *App {
	logger := logger.NewLogger()

	commonHandler := handlers.NewCommonHandler(logger)

	bookService := book.NewBookService(db)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &App{
		BookHandler: *bookHandler,
	}
	return app
}
