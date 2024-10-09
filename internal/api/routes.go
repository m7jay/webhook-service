package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/m7jay/webhook-service/internal/api/handlers"
	"github.com/m7jay/webhook-service/internal/api/middlewares"
	"github.com/m7jay/webhook-service/internal/services"
	"github.com/m7jay/webhook-service/internal/utils"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, redisClient *utils.RedisClient, kafkaProducer *utils.KafkaProducer) {
	webhookService := services.NewWebhookService(db, redisClient, kafkaProducer)
	eventHandler := handlers.NewEventHandler(db, webhookService)
	subscriptionHandler := handlers.NewSubscriptionHandler(db)

	api := app.Group("/api") //, middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")))

	// Event routes
	events := api.Group("/events")
	events.Post("/", middlewares.RBACMiddleware("admin"), eventHandler.CreateEvent)
	// events.Get("/", eventHandler.ListEvents)
	// events.Get("/:id", eventHandler.GetEvent)
	// events.Put("/:id", middlewares.RBACMiddleware("admin"), eventHandler.UpdateEvent)
	// events.Delete("/:id", middlewares.RBACMiddleware("admin"), eventHandler.DeleteEvent)
	events.Post("/:id/trigger", middlewares.RBACMiddleware("admin"), eventHandler.TriggerManualEvent)
	events.Post("/:id/upload", middlewares.RBACMiddleware("admin"), eventHandler.UploadFile)

	// Subscription routes
	subscriptions := api.Group("/subscriptions")
	subscriptions.Post("/", subscriptionHandler.CreateSubscription)
	// subscriptions.Get("/", subscriptionHandler.ListSubscriptions)
	// subscriptions.Get("/:id", subscriptionHandler.GetSubscription)
	// subscriptions.Put("/:id", subscriptionHandler.UpdateSubscription)
	// subscriptions.Delete("/:id", subscriptionHandler.DeleteSubscription)

	// Metrics endpoint
	app.Get("/metrics", middlewares.MetricsMiddleware())
	app.Get("/api/check", handlers.HealthCheck)
}
