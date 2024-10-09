// utils/kafka.go
package utils

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/m7jay/webhook-service/config"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func InitKafka(cfg *config.Config) (*KafkaProducer, *kafka.Reader) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.KafkaTopic,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.KafkaTopic,
		GroupID: cfg.KafkaGroupID,
	})

	log.Println("Kafka producer and consumer initialized successfully")
	return &KafkaProducer{Writer: writer}, reader
}

func (kp *KafkaProducer) Produce(key, value []byte) error {
	return kp.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)
}
