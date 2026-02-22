package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) accessLogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startedAt := time.Now()

		err := c.Next()

		latency := time.Since(startedAt)
		requestID, _ := c.Locals("requestid").(string)
		s.logger.Info("http_request",
			"request_id", requestID,
			"status", c.Response().StatusCode(),
			"method", c.Method(),
			"path", c.OriginalURL(),
			"ip", c.IP(),
			"latency_ms", latency.Milliseconds(),
		)

		return err
	}
}
