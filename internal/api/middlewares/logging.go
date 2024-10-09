// logging.go
package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		logger.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"latency":    c.Response().Header.Peek("X-Response-Time"),
			"ip":         c.IP(),
			"user-agent": c.Get("User-Agent"),
		}).Info("Request processed")

		return err
	}
}
