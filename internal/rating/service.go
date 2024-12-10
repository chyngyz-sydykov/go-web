package rating

import (
	"context"
	"time"

	my_error "github.com/chyngyz-sydykov/go-web/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/chyngyz-sydykov/go-web/proto/rating"
)

type RatingServiceInterface interface {
	GetByBookId(bookId int) ([]RatingDTO, error)
	Create(ratingDTO *RatingDTO) error
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

	ratingsDto := service.mapProtoBuffRatingListToRatingDTOList(response.Ratings)
	return ratingsDto, nil
}

func (service *RatingService) Create(ratingDTO *RatingDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), service.timeout)
	defer cancel()
	saveRatingRequest := service.maRatingDTOTopProtoBuffRating(ratingDTO)
	response, err := service.client.SaveRating(ctx, saveRatingRequest)

	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return my_error.ErrInvalidArgument
		}
		return my_error.ErrgRpcServerDown

	}

	service.mapProtoBuffRatingToRatingDTO(ratingDTO, response.Rating)
	return nil
}

func (s *RatingService) mapProtoBuffRatingToRatingDTO(ratingDTO *RatingDTO, ratings *pb.Rating) {
	ratingDTO.RatingID = ratings.GetRatingId()
	ratingDTO.BookID = ratings.GetBookId()
	ratingDTO.Rating = ratings.GetRating()
	ratingDTO.Comment = ratings.GetComment()

}
func (s *RatingService) maRatingDTOTopProtoBuffRating(ratingDTO *RatingDTO) *pb.SaveRatingRequest {
	return &pb.SaveRatingRequest{
		BookId:  ratingDTO.BookID,
		Rating:  ratingDTO.Rating,
		Comment: ratingDTO.Comment,
	}
}

func (s *RatingService) mapProtoBuffRatingListToRatingDTOList(ratings []*pb.Rating) []RatingDTO {
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
