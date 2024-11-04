package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string
	ICBN        string
	PublishedAt *time.Time
	AuthorId    int8
	Author      Author
}
