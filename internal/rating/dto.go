package rating

type RatingDTO struct {
	RatingID string `json:"rating_id"`
	BookID   int32  `json:"book_id"`
	Rating   int32  `json:"rating"`
	Comment  string `json:"comment,omitempty"` // Optional field
}
