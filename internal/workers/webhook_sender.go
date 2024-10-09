package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/m7jay/webhook-service/config"
	"github.com/m7jay/webhook-service/internal/services"
	"github.com/m7jay/webhook-service/internal/utils"
)

type WebhookSender struct {
	redisClient *utils.RedisClient
	service     *services.WebhookService
	queueName   string
}

func NewWebhookSender(cfg *config.Config, redisClient *utils.RedisClient, service *services.WebhookService) *WebhookSender {
	return &WebhookSender{
		redisClient: redisClient,
		service:     service,
		queueName:   "webhook_queue",
	}
}

func (ws *WebhookSender) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go ws.worker()
	}
}

func (ws *WebhookSender) worker() {
	ctx := context.Background()
	for {
		result, err := ws.redisClient.Client.BLPop(ctx, 0, ws.queueName).Result()
		if err != nil {
			log.Printf("Error popping from queue: %v", err)
			continue
		}

		// The second element is the value
		jobJSON := result[1]

		var job struct {
			SubscriptionID uint        `json:"subscription_id"`
			Payload        interface{} `json:"payload"`
		}

		if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
			log.Printf("Error unmarshalling job: %v", err)
			continue
		}

		if err := ws.service.SendWebhook(job.SubscriptionID, job.Payload); err != nil {
			log.Printf("Error sending webhook: %v", err)
			// Implement retry logic here
			ws.retryJob(job)
		}
	}
}

func (ws *WebhookSender) EnqueueWebhook(subscriptionID uint, payload interface{}) error {
	job := struct {
		SubscriptionID uint        `json:"subscription_id"`
		Payload        interface{} `json:"payload"`
	}{
		SubscriptionID: subscriptionID,
		Payload:        payload,
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("error marshalling job: %w", err)
	}

	ctx := context.Background()
	return ws.redisClient.Client.RPush(ctx, ws.queueName, jobJSON).Err()
}

func (ws *WebhookSender) retryJob(job struct {
	SubscriptionID uint        `json:"subscription_id"`
	Payload        interface{} `json:"payload"`
}) {
	// Implement exponential backoff retry logic here
	// For simplicity, we'll just re-enqueue the job with a delay
	time.Sleep(5 * time.Second)
	if err := ws.EnqueueWebhook(job.SubscriptionID, job.Payload); err != nil {
		log.Printf("Error re-enqueueing job: %v", err)
	}
}
