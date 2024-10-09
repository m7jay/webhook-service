package models

import (
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	UserID   uint
	EventID  uint
	Endpoint string
	Secret   string
	IsActive bool
	Event    Event `gorm:"foreignKey:EventID"`
}
