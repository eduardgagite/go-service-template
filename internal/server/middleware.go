package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// authMiddleware — точка расширения для аутентификации. Сейчас пропускает все
// запросы, чтобы шаблон работал из коробки. Замените тело на реальную
// аутентификацию (JWT / API key / OIDC): проверьте учётные данные и отклоняйте
// неаутентифицированные запросы, например:
//
//	token := c.Get(fiber.HeaderAuthorization)
//	if !valid(token) {
//	    return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "unauthorized"})
//	}
func (s *Server) authMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: аутентифицировать запрос до попадания в обработчики.
		return c.Next()
	}
}

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
