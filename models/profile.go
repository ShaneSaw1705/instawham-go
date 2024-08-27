package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Username string
	Mood     string
	AuthID   int `gorm:"unique"`
}
