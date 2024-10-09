// utils/redis.go
package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gocraft/work"

	"github.com/m7jay/webhook-service/config"
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedis(cfg *config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
	return &RedisClient{Client: client}
}

func (rc *RedisClient) Enqueue(job *work.Job) error {
	// Implement job enqueuing logic here
	// This is a placeholder and needs to be implemented based on your specific requirements
	return nil
}
