package rating

import (
	"context"
	"log"
	"time"

	my_error "github.com/chyngyz-sydykov/go-web/error"

	pb "github.com/chyngyz-sydykov/go-web/proto/rating"
)

type RatingServiceInterface interface {
	GetByBookId(bookId uint) ([]*pb.Rating, error)
	//Create(rating any) error
}

type RatingService struct {
	client  pb.RatingServiceClient
	timeout time.Duration
}

func NewRatingService(client pb.RatingServiceClient, timeout time.Duration) *RatingService {

	return &RatingService{client: client, timeout: timeout}
}

func (service *RatingService) GetByBookId(bookId uint) ([]*pb.Rating, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.timeout)
	defer cancel()

	getRatingRequest := &pb.GetRatingsRequest{
		BookId: int32(bookId),
	}
	response, err := service.client.GetRatings(ctx, getRatingRequest)
	if err != nil {
		return nil, my_error.ErrgRpcServerDown

	}

	log.Printf("Rating returned successfully: %v", response.Ratings)
	return response.Ratings, nil
}

// func (service *RatingService) Create(rating any) error {
// 	//return service.repository.Create(book)
// 	return my_error.ErrNotFound
// }
