package book

import (
	"fmt"
	"log"
	"time"

	my_error "github.com/chyngyz-sydykov/go-web/error"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	messagebroker "github.com/chyngyz-sydykov/go-web/infrastructure/messagebroker"
	"github.com/chyngyz-sydykov/go-web/internal/rating"

	"gorm.io/gorm"
)

const BOOKUPDATED_EVENTNAME = "bookUpdated"

type BookServiceInterface interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (BookDTO, error)
	Create(book BookRequest) (*models.Book, error)
	Update(id int, payload BookRequest) (*models.Book, error)
	Delete(id int) error
}

type BookService struct {
	repository    BookRepository
	messageBroker messagebroker.MessageBrokerInterface
	ratingService rating.RatingServiceInterface
}

func NewBookService(db *gorm.DB, messageBrokerPublisher messagebroker.MessageBrokerInterface, ratingService rating.RatingServiceInterface) *BookService {
	repository := NewBookRepository(db)
	return &BookService{
		repository:    *repository,
		messageBroker: messageBrokerPublisher,
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

func (service *BookService) Create(bookRequest BookRequest) (*models.Book, error) {
	gormBook := service.mapBookRequestToGormBook(bookRequest)
	publishedAt := time.Now()
	gormBook.PublishedAt = &publishedAt
	return service.repository.Create(gormBook)
}

func (service *BookService) Update(id int, bookRequest BookRequest) (*models.Book, error) {
	var book models.Book

	// Find the book by ID
	book, err := service.repository.GetByID(id)
	if err != nil {
		return nil, my_error.ErrNotFound
	}

	// Update the book with the payload fields
	payload := service.mapBookRequestToGormBook(bookRequest)
	err = service.repository.Update(&book, *payload)
	if err != nil {
		return nil, err
	}

	err = service.publishMessage(book, BOOKUPDATED_EVENTNAME)
	if err != nil {
		fmt.Println("Failed to publish message")
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

func (service *BookService) mapBookRequestToGormBook(bookRequest BookRequest) *models.Book {
	return &models.Book{
		Title:    bookRequest.Title,
		ICBN:     bookRequest.ICBN,
		AuthorId: int64(bookRequest.AuthorID),
	}
}

func (service *BookService) publishMessage(book models.Book, event string) error {

	bookMessage := BookMessage{
		BookId:   int(book.ID),
		EditedAt: time.Now(),
		Event:    event,
	}

	if err := service.messageBroker.Publish(bookMessage); err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}
	return nil
}
