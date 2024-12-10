package rating

type RatingDTO struct {
	RatingID string `json:"ratingId"`
	BookID   int32  `json:"bookId"`
	Rating   int32  `json:"rating"`
	Comment  string `json:"comment,omitempty"` // Optional field
}
