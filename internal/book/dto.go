package book

import (
	"time"

	"github.com/chyngyz-sydykov/go-web/internal/rating"
)

type AuthorDTO struct {
	ID        uint
	Firstname string
	Lastname  string
}

type BookDTO struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	ICBN        string             `json:"icbn"`
	PublishedAt *time.Time         `json:"published_at"`
	AuthorID    int32              `json:"author_id"`
	Author      AuthorDTO          `json:"author"`
	Ratings     []rating.RatingDTO `json:"ratings"`
}

type BookMessage struct {
	BookId   int       `json:"book_id"`
	EditedAt time.Time `json:"EditedAt"`
	Event    string    `json:"event"`
}
