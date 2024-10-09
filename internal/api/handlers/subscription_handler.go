package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m7jay/webhook-service/internal/models"
	"github.com/m7jay/webhook-service/internal/services"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	DB      *gorm.DB
	Service *services.SubscriptionService
}

func NewSubscriptionHandler(db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{
		DB:      db,
		Service: services.NewSubscriptionService(db),
	}
}

func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	subscription := new(models.Subscription)
	if err := c.BodyParser(subscription); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.Service.CreateSubscription(subscription); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create subscription"})
	}

	return c.Status(fiber.StatusCreated).JSON(subscription)
}

// Implement other subscription-related handlers (GetSubscription, ListSubscriptions, UpdateSubscription, DeleteSubscription)
