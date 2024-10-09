package services

import (
	"errors"

	"github.com/m7jay/webhook-service/internal/models"
	"gorm.io/gorm"
)

type SubscriptionService struct {
	db *gorm.DB
}

func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

func (s *SubscriptionService) CreateSubscription(subscription *models.Subscription) error {
	return s.db.Create(subscription).Error
}

func (s *SubscriptionService) GetSubscription(id uint) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := s.db.First(&subscription, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}
	return &subscription, nil
}

func (s *SubscriptionService) UpdateSubscription(subscription *models.Subscription) error {
	return s.db.Save(subscription).Error
}

func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.db.Delete(&models.Subscription{}, id).Error
}

func (s *SubscriptionService) ListSubscriptions(page, pageSize int) ([]models.Subscription, int64, error) {
	var subscriptions []models.Subscription
	var total int64

	offset := (page - 1) * pageSize

	if err := s.db.Model(&models.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset(offset).Limit(pageSize).Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (s *SubscriptionService) GetSubscriptionsByEvent(eventID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := s.db.Where("event_id = ?", eventID).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *SubscriptionService) GetActiveSubscriptionsByEvent(eventID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := s.db.Where("event_id = ? AND is_active = ?", eventID, true).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}
