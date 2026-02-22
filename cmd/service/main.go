package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go-service-template/internal/config"
	"go-service-template/internal/server"
	"go-service-template/internal/service"
	"go-service-template/internal/storage/postgres"
)

// @title Service API
// @version 1.0
// @description REST API for microservice template
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @tag.name health
// @tag.description Health check endpoints

// @tag.name examples
// @tag.description Example CRUD operations

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "service failed to start:", err)
		os.Exit(1)
	}
}

func run() error {
	app, err := NewApp()
	if err != nil {
		return err
	}

	return app.Run()
}

type App struct {
	cfg     *config.Config
	logger  *slog.Logger
	logFile *os.File
	db      *postgres.PostgresStorage
	server  *server.Server
}

func NewApp() (*App, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	logFile, logger, err := setupLogger(cfg.App.DebugMode)
	if err != nil {
		return nil, err
	}

	db, err := initStorage(cfg)
	if err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("connect database: %w", err)
	}

	services := service.NewServices(db, logger)
	srv := server.New(services, logger, cfg)

	return &App{
		cfg:     cfg,
		logger:  logger,
		logFile: logFile,
		db:      db,
		server:  srv,
	}, nil
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return cfg, nil
}

func initStorage(cfg *config.Config) (*postgres.PostgresStorage, error) {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	db, err := postgres.NewStorage(dbCtx, cfg.DatabaseDSN(), cfg.Database)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverErr := make(chan error, 1)

	go func() {
		if err := a.Start(); err != nil {
			serverErr <- fmt.Errorf("server run: %w", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case <-quit:
		a.logger.Info("Shutdown signal received")
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		a.logger.Info("Context cancelled")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	return a.Shutdown(shutdownCtx)
}

func (a *App) Start() error {
	portStr := strconv.Itoa(a.cfg.Server.Port)
	return a.server.Start(portStr)
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Server is shutting down...")

	var shutdownErrs []error

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		shutdownErrs = append(shutdownErrs, fmt.Errorf("shutdown server: %w", err))
	} else {
		a.logger.Info("Server exited gracefully")
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error("Failed to close database", slog.String("error", err.Error()))
			shutdownErrs = append(shutdownErrs, fmt.Errorf("close database: %w", err))
		}
	}

	if a.logFile != nil {
		if err := a.logFile.Close(); err != nil {
			shutdownErrs = append(shutdownErrs, fmt.Errorf("close log file: %w", err))
		}
	}

	return errors.Join(shutdownErrs...)
}

func setupLogger(debugMode bool) (*os.File, *slog.Logger, error) {
	var logger *slog.Logger
	logFile, err := createLogFile()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "log file is unavailable, using stdout only: %v\n", err)
		if debugMode {
			logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		} else {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}
		return nil, logger, nil
	}

	if debugMode {
		logger = slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logFile, logger, nil
}

func createLogFile() (*os.File, error) {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	logFilePath := "logs/" + "logfile_" + time.Now().Format("02-01-2006_15-04-05") + ".log"

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("open log file %q: %w", logFilePath, err)
	}

	return logFile, nil
}
