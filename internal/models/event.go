package models

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Description string
}
