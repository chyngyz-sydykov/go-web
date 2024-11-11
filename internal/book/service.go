package book

import (
	"fmt"

	"github.com/chyngyz-sydykov/go-web/db/models"

	"gorm.io/gorm"
)

type BookServiceInterface interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (models.Book, error)
	Create(book *models.Book) error
	Update(book *models.Book) error
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

func (service *BookService) GetByID(id uint) (models.Book, error) {
	return service.repository.GetByID(id)
}

func (service *BookService) Create(book *models.Book) error {
	return service.repository.Create(book)
}

func (service *BookService) Update(id uint, payload models.Book) (*models.Book, error) {
	var book models.Book

	// Find the book by ID
	fmt.Println("salam1")
	book, err := service.repository.GetByID(id)
	if err != nil {
		fmt.Println("salam2")
		return nil, err
	}

	fmt.Println("salam3")

	// Update the book with the payload fields
	err = service.repository.Update(&book, payload)
	if err != nil {
		return nil, err
	}

	return &book, nil
}
