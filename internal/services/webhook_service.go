package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocraft/work"
	"gorm.io/gorm"

	"github.com/m7jay/webhook-service/internal/utils"

	"github.com/m7jay/webhook-service/internal/models"
)

type WebhookService struct {
	DB            *gorm.DB
	RedisClient   *utils.RedisClient
	KafkaProducer *utils.KafkaProducer
}

func NewWebhookService(db *gorm.DB, redisClient *utils.RedisClient, kafkaProducer *utils.KafkaProducer) *WebhookService {
	return &WebhookService{
		DB:            db,
		RedisClient:   redisClient,
		KafkaProducer: kafkaProducer,
	}
}

func (s *WebhookService) TriggerWebhook(eventID uint, payload interface{}) error {
	var subscriptions []models.Subscription
	if err := s.DB.Where("event_id = ? AND is_active = ?", eventID, true).Find(&subscriptions).Error; err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		// Enqueue webhook job
		job := &work.Job{
			Name: "send_webhook",
			ID:   utils.GenerateUniqueID(),
			Args: work.Q{
				"subscription_id": subscription.ID,
				"payload":         payload,
			},
		}

		if err := s.RedisClient.Enqueue(job); err != nil {
			return err
		}
	}

	return nil
}

func (s *WebhookService) SendWebhook(subscriptionID uint, payload interface{}) error {
	var subscription models.Subscription
	if err := s.DB.First(&subscription, subscriptionID).Error; err != nil {
		return err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send HTTP request
	resp, err := http.Post(subscription.Endpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return s.handleWebhookFailure(subscription, payload, err)
	}
	defer resp.Body.Close()

	// Log webhook attempt
	log := models.WebhookLog{
		EventID:        subscription.EventID,
		SubscriptionID: subscription.ID,
		Payload:        string(payloadBytes),
		ResponseCode:   resp.StatusCode,
		Success:        resp.StatusCode >= 200 && resp.StatusCode < 300,
	}
	s.DB.Create(&log)

	if log.Success {
		return nil
	}

	return s.handleWebhookFailure(subscription, payload, nil)
}

func (s *WebhookService) handleWebhookFailure(subscription models.Subscription, payload interface{}, err error) error {
	// Implement exponential backoff retry logic
	retryCount := 0
	maxRetries := 5
	baseDelay := time.Second

	for retryCount < maxRetries {
		delay := time.Duration(1<<uint(retryCount)) * baseDelay
		time.Sleep(delay)

		if err := s.SendWebhook(subscription.ID, payload); err == nil {
			return nil
		}

		retryCount++
	}

	// If all retries fail, log the final failure
	failureLog := models.WebhookLog{
		EventID:        subscription.EventID,
		SubscriptionID: subscription.ID,
		Payload:        payload.(string),
		ResponseCode:   0,
		ResponseBody:   err.Error(),
		Attempts:       maxRetries,
		Success:        false,
	}
	s.DB.Create(&failureLog)

	return err
}
