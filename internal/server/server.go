package server

import (
	"context"
	"log/slog"
	"net"
	"time"

	"go-service-template/internal/config"
	"go-service-template/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
		BodyLimit:    s.config.Server.BodyLimit,
	})

	s.app.Use(recover.New())
	s.app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
	}))
	s.app.Use(helmet.New())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: s.config.Server.CORSAllowOrigins,
	}))
	if s.config.Server.RateLimit > 0 {
		s.app.Use(limiter.New(limiter.Config{
			Max:        s.config.Server.RateLimit,
			Expiration: time.Minute,
		}))
	}
	s.app.Use(s.accessLogMiddleware())

	// Swagger раскрывает всю поверхность API, поэтому закрыт флагом ENABLE_SWAGGER
	// (по умолчанию выключен). Каталог docs/ генерируется на этапе сборки.
	if s.config.App.EnableSwagger {
		s.app.Static("/swagger/docs", "./docs")
		s.app.Get("/swagger/*", swagger.New(swagger.Config{
			URL: "/swagger/docs/swagger.json",
		}))
	}

	// Пробы: /livez без зависимостей (для k8s livenessProbe); /readyz пингует
	// базу (для k8s readinessProbe). /health оставлен как обратно совместимый
	// алиас readiness.
	s.app.Get("/livez", s.liveness)
	s.app.Get("/readyz", s.readiness)
	s.app.Get("/health", s.readiness)

	api := s.app.Group("/api/v1")
	// authMiddleware пока пропускает все запросы — замените на реальную аутентификацию.
	api.Use(s.authMiddleware())

	examples := api.Group("/examples")
	examples.Post("/", s.createExample)
	examples.Get("/", s.getAllExamples)
	examples.Get("/:id", s.getExample)
	examples.Put("/:id", s.updateExample)
	examples.Delete("/:id", s.deleteExample)
}

func (s *Server) Start(port string) error {
	s.setupRoutes()

	addr := net.JoinHostPort(s.config.Server.Host, port)
	s.logger.Info("Starting server",
		slog.String("host", s.config.Server.Host),
		slog.String("port", port),
		slog.String("addr", addr),
	)

	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")

	if s.app == nil {
		return nil
	}

	return s.app.ShutdownWithContext(ctx)
}
