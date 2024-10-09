// errors.go
package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return nil
	}
}
