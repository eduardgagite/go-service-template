package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
    MaxConns int
    MinConns int
    MaxConnLifetime time.Duration
    MaxConnIdleTime time.Duration
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppConfig struct {
	DebugMode bool
}

func Load() (*Config, error) {
	config := &Config{}

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnvInt("DB_PORT", 5432)
	config.Database.Name = getEnv("DB_NAME", "service_db")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
    config.Database.MaxConns = getEnvInt("DB_MAX_CONNS", 10)
    config.Database.MinConns = getEnvInt("DB_MIN_CONNS", 1)
    config.Database.MaxConnLifetime = getEnvDuration("DB_MAX_CONN_LIFETIME", time.Hour)
    config.Database.MaxConnIdleTime = getEnvDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute)

	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnvInt("SERVER_PORT", 8080)
	config.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	config.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)

	config.App.DebugMode = getEnvBool("DEBUG_MODE", false)

	return config, nil
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
