package server

import (
	"go-service-template/internal/storage"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

type Server struct {
	storage storage.Storage
	logger  *slog.Logger
	app     *fiber.App
}

func New(storage storage.Storage, slogger *slog.Logger) *Server {
	return &Server{
		storage: storage,
		logger:  slogger,
	}
}

func (s *Server) setupRoutes() {
	s.app = fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
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
