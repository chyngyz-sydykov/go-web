package rating

import (
	"context"
	"log"
	"time"

	my_error "github.com/chyngyz-sydykov/go-web/error"

	pb "github.com/chyngyz-sydykov/go-web/proto/rating"
)

type RatingServiceInterface interface {
	GetByBookId(bookId int) ([]RatingDTO, error)
	//Create(rating any) error
}

type RatingService struct {
	client  pb.RatingServiceClient
	timeout time.Duration
}

func NewRatingService(client pb.RatingServiceClient, timeout time.Duration) *RatingService {

	return &RatingService{client: client, timeout: timeout}
}

func (service *RatingService) GetByBookId(bookId int) ([]RatingDTO, error) {
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
	ratingDto := service.mapProtoBuffRatingToRatingDTO(response.Ratings)
	return ratingDto, nil
}

// func (service *RatingService) Create(rating any) error {
// 	//return service.repository.Create(book)
// 	return my_error.ErrNotFound
// }

func (s *RatingService) mapProtoBuffRatingToRatingDTO(ratings []*pb.Rating) []RatingDTO {
	mapped := make([]RatingDTO, len(ratings))

	for i, r := range ratings {
		mapped[i] = RatingDTO{
			RatingID: r.GetRatingId(),
			BookID:   r.GetBookId(),
			Rating:   r.GetRating(),
			Comment:  r.GetComment(),
		}
	}

	return mapped
}
