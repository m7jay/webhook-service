package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// type HealthCheckHandler struct {
// 	status int
// }

// func NewHealthCheckHandler() *HealthCheckHandler {
// 	return &HealthCheckHandler{
// 		status: fiber.StatusOK,
// 	}
// }

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success!"})
}
