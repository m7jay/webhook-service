package models

import (
	"gorm.io/gorm"
)

type WebhookLog struct {
	gorm.Model
	EventID        uint
	SubscriptionID uint
	Payload        string
	ResponseCode   int
	ResponseBody   string
	Attempts       int
	Success        bool
}
