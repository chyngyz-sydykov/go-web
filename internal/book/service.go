package book

import (
	my_error "github.com/chyngyz-sydykov/go-web/error"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-web/internal/rating"

	"gorm.io/gorm"
)

type BookServiceInterface interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (BookDTO, error)
	Create(book *models.Book) error
	Update(id int, payload models.Book) (*models.Book, error)
	Delete(id int) error
}

type BookService struct {
	repository    BookRepository
	ratingService rating.RatingServiceInterface
}

func NewBookService(db *gorm.DB, ratingService rating.RatingServiceInterface) *BookService {
	repository := NewBookRepository(db)
	return &BookService{
		repository:    *repository,
		ratingService: ratingService,
	}
}

func (service *BookService) GetAll() ([]models.Book, error) {
	return service.repository.GetAll()
}

func (service *BookService) GetByID(id int) (BookDTO, error) {
	book, err := service.repository.GetByID(id)
	if err != nil {
		return BookDTO{}, err
	}
	ratings, err := service.ratingService.GetByBookId(id)
	if err != nil {
		return service.mapToBookingDTO(book, ratings), err
	}
	return service.mapToBookingDTO(book, ratings), nil
}

func (service *BookService) Create(book *models.Book) error {
	return service.repository.Create(book)
}

func (service *BookService) Update(id int, payload models.Book) (*models.Book, error) {
	var book models.Book

	// Find the book by ID
	book, err := service.repository.GetByID(id)
	if err != nil {
		return nil, my_error.ErrNotFound
	}

	// Update the book with the payload fields
	err = service.repository.Update(&book, payload)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (service *BookService) Delete(id int) error {
	var book models.Book

	// Find the book by ID
	book, err := service.repository.GetByID(id)
	if err != nil {
		return my_error.ErrNotFound
	}

	// Update the book with the payload fields
	err = service.repository.Delete(&book)
	if err != nil {
		return err
	}

	return nil
}

func (service *BookService) mapToBookingDTO(book models.Book, ratingDTO []rating.RatingDTO) BookDTO {
	bookDTO := BookDTO{
		ID:          book.ID,
		Title:       book.Title,
		ICBN:        book.ICBN,
		PublishedAt: book.PublishedAt,
		Ratings:     ratingDTO,
		AuthorID:    int32(book.Author.ID),
		Author: AuthorDTO{
			ID:        book.Author.ID,
			Firstname: book.Author.Firstname,
			Lastname:  book.Author.Lastname,
		},
	}
	return bookDTO
}
