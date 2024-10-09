package workers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/m7jay/webhook-service/internal/services"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	service *services.WebhookService
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, service *services.WebhookService) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		service: service,
	}
}

func (c *KafkaConsumer) Start() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var event struct {
			EventID uint        `json:"event_id"`
			Payload interface{} `json:"payload"`
		}

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if err := c.service.TriggerWebhook(event.EventID, event.Payload); err != nil {
			log.Printf("Error triggering webhook: %v", err)
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
