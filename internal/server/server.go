package server

import (
	"context"
	"log/slog"
	"time"

	"go-service-template/internal/config"
	"go-service-template/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

type Server struct {
	services *service.Services
	logger   *slog.Logger
	config   *config.Config
	app      *fiber.App
}

func New(services *service.Services, slogger *slog.Logger, cfg *config.Config) *Server {
	return &Server{
		services: services,
		logger:   slogger,
		config:   cfg,
	}
}

func (s *Server) setupRoutes() {
	s.app = fiber.New(fiber.Config{
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	})

	s.app.Use(recover.New())
	s.app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	s.app.Get("/swagger/*", swagger.HandlerDefault)

	s.app.Get("/health", s.healthCheck)

	api := s.app.Group("/api/v1")

	examples := api.Group("/examples")
	examples.Post("/", s.createExample)
	examples.Get("/", s.getAllExamples)
	examples.Get("/:id", s.getExample)
	examples.Put("/:id", s.updateExample)
	examples.Delete("/:id", s.deleteExample)
}

func (s *Server) Start(port string) error {
	s.setupRoutes()

	s.logger.Info("Starting server", slog.String("port", port))

	return s.app.Listen(":" + port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")

	if s.app == nil {
		return nil
	}

	return s.app.ShutdownWithContext(ctx)
}
