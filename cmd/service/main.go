package main

import (
	"context"
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
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	moscowLoc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		moscowLoc = time.FixedZone("MSK", 3*60*60)
	}
	time.Local = moscowLoc

	logFile, logger := setupLogger(cfg.App.DebugMode)
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	db, err := postgres.NewStorage(cfg.DatabaseDSN())
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	services := service.NewServices(db, logger)
	srv := server.New(services, logger, cfg)
	port := strconv.Itoa(cfg.Server.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info("Starting server", slog.String("port", port))
		if err := srv.Start(port); err != nil {
			logger.Error("Server failed to start", slog.String("error", err.Error()))
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("Shutdown signal received")
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	logger.Info("Server is shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	} else {
		logger.Info("Server exited gracefully")
	}
}

func setupLogger(debugMode bool) (*os.File, *slog.Logger) {
	var logger *slog.Logger
	logFile := createLogFile()

	if debugMode {
		logger = slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logFile, logger
}

func createLogFile() *os.File {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	logFilePath := "logs/" + "logfile_" + time.Now().Format("02-01-2006_15-04-05") + ".log"

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	return logFile
}
