package services

import (
	"errors"

	"github.com/m7jay/webhook-service/internal/models"
	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	return s.db.Create(event).Error
}

func (s *EventService) GetEvent(id uint) (*models.Event, error) {
	var event models.Event
	if err := s.db.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	return s.db.Save(event).Error
}

func (s *EventService) DeleteEvent(id uint) error {
	return s.db.Delete(&models.Event{}, id).Error
}

func (s *EventService) ListEvents(page, pageSize int) ([]models.Event, int64, error) {
	var events []models.Event
	var total int64

	offset := (page - 1) * pageSize

	if err := s.db.Model(&models.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
