package book

import (
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"

	"gorm.io/gorm"
)

type BookRepositoryInterface interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (models.Book, error)
	Create(book models.Book) error
	Update(book *models.Book, payload models.Book) error
}

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (repository *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	if err := repository.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (repository *BookRepository) GetByID(id uint) (models.Book, error) {
	var book models.Book
	if err := repository.db.First(&book, id).Error; err != nil {
		return book, err
	}
	return book, nil
}

func (repository *BookRepository) Create(book *models.Book) error {
	return repository.db.Create(&book).Error
}

func (repository *BookRepository) Update(book *models.Book, payload models.Book) error {
	if err := repository.db.Model(&book).Updates(payload).Error; err != nil {
		return err
	}
	return nil
}

func (repository *BookRepository) Delete(book *models.Book) error {
	if err := repository.db.Unscoped().Delete(&models.Book{}, book.ID).Error; err != nil {
		return err
	}
	return nil
}
