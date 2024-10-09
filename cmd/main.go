package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/m7jay/webhook-service/config"
	routes "github.com/m7jay/webhook-service/internal/api"
	"github.com/m7jay/webhook-service/internal/api/middlewares"
	"github.com/m7jay/webhook-service/internal/models"
	"github.com/m7jay/webhook-service/internal/utils"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize database
	db, err := utils.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Event{}, &models.Subscription{}, &models.WebhookLog{}); err != nil {
		log.Fatal("Failed to auto-migrate models:", err)
	}

	// Initialize Redis
	redisClient := utils.InitRedis(cfg)

	// Initialize Kafka
	kafkaProducer, _ := utils.InitKafka(cfg)

	// Initialize Fiber app
	app := fiber.New()

	// Apply global middlewares
	app.Use(middlewares.MetricsMiddleware())
	app.Use(middlewares.LoggingMiddleware(cfg.Logger))
	app.Use(middlewares.ErrorMiddleware())
	// app.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	// Setup routes
	routes.SetupRoutes(app, db, redisClient, kafkaProducer)

	// Start the server
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
