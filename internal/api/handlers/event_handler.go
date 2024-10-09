// Implement other event-related handlers (GetEvent, ListEvents, UpdateEvent, DeleteEvent)
package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/m7jay/webhook-service/internal/models"
	"gorm.io/gorm"

	"github.com/m7jay/webhook-service/internal/services"
	"github.com/m7jay/webhook-service/internal/workers"
)

type EventHandler struct {
	DB             *gorm.DB
	Service        *services.EventService
	WebhookService *services.WebhookService
	FileProcessor  *workers.FileProcessor
}

func NewEventHandler(db *gorm.DB, webhookService *services.WebhookService) *EventHandler {
	return &EventHandler{
		DB:             db,
		Service:        services.NewEventService(db),
		WebhookService: webhookService,
		FileProcessor:  workers.NewFileProcessor(webhookService),
	}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	event := new(models.Event)
	if err := c.BodyParser(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.Service.CreateEvent(event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create event"})
	}

	return c.Status(fiber.StatusCreated).JSON(event)
}

func (h *EventHandler) TriggerManualEvent(c *fiber.Ctx) error {
	eventID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.WebhookService.TriggerWebhook(uint(eventID), payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to trigger webhook"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Event triggered successfully"})
}

func (h *EventHandler) UploadFile(c *fiber.Ctx) error {
	eventID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get uploaded file"})
	}

	uploadDir := "./uploads"
	filename := filepath.Join(uploadDir, file.Filename)

	if err := c.SaveFile(file, filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	go func() {
		if err := h.FileProcessor.ProcessFile(filename, uint(eventID)); err != nil {
			fmt.Printf("Error processing file: %v\n", err)
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File uploaded and processing started"})
}
