package book

import (
	"github.com/chyngyz-sydykov/go-web/db/models"

	"gorm.io/gorm"
)

type BookServiceInterface interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (models.Book, error)
	Create(book models.Book) error
}

type BookService struct {
	repository BookRepository
}

func NewBookService(db *gorm.DB) *BookService {
	repository := NewBookRepository(db)
	return &BookService{repository: *repository}
}

func (service *BookService) GetAll() ([]models.Book, error) {
	return service.repository.GetAll()
}

func (service *BookService) GetByID(id int) (models.Book, error) {
	return service.repository.GetByID(id)
}

func (service *BookService) Create(book models.Book) error {
	return service.repository.Create(book)
}
