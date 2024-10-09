// metrics.go
package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		httpRequestsTotal.WithLabelValues(c.Method(), c.Path(), string(c.Response().Header.Peek("Status"))).Inc()
		return err
	}
}
