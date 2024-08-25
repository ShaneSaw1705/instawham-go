package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string
	Description string
	AuthorID    int
}
