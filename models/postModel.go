package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	title       string
	description string
	authorID    int
}
