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
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type ServerConfig struct {
	Host             string
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	BodyLimit        int    // максимальный размер тела запроса в байтах
	RateLimit        int    // лимит запросов в минуту на один IP (0 отключает лимитер)
	CORSAllowOrigins string // список разрешённых CORS-источников через запятую
}

type AppConfig struct {
	DebugMode bool
	// EnableSwagger включает эндпоинты Swagger UI / docs. В продакшене держите
	// выключенным: они раскрывают всю поверхность API. По умолчанию false.
	EnableSwagger bool
}

func Load() (*Config, error) {
	config := &Config{}
	var err error

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port, err = getEnvInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}
	config.Database.Name = getEnv("DB_NAME", "service_db")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	// "disable" подходит для локальной сети docker-compose. Для любого удалённого
	// или managed-Postgres в продакшене используйте "require" (или "verify-full" с CA).
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	config.Database.MaxConns, err = getEnvInt("DB_MAX_CONNS", 10)
	if err != nil {
		return nil, err
	}
	config.Database.MinConns, err = getEnvInt("DB_MIN_CONNS", 1)
	if err != nil {
		return nil, err
	}
	config.Database.MaxConnLifetime, err = getEnvDuration("DB_MAX_CONN_LIFETIME", time.Hour)
	if err != nil {
		return nil, err
	}
	config.Database.MaxConnIdleTime, err = getEnvDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute)
	if err != nil {
		return nil, err
	}

	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port, err = getEnvInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, err
	}
	config.Server.ReadTimeout, err = getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}
	config.Server.WriteTimeout, err = getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return nil, err
	}
	config.Server.BodyLimit, err = getEnvInt("SERVER_BODY_LIMIT", 4*1024*1024)
	if err != nil {
		return nil, err
	}
	config.Server.RateLimit, err = getEnvInt("SERVER_RATE_LIMIT", 100)
	if err != nil {
		return nil, err
	}
	config.Server.CORSAllowOrigins = getEnv("CORS_ALLOW_ORIGINS", "*")

	config.App.DebugMode, err = getEnvBool("DEBUG_MODE", false)
	if err != nil {
		return nil, err
	}
	config.App.EnableSwagger, err = getEnvBool("ENABLE_SWAGGER", false)
	if err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Database.Name == "" {
		return fmt.Errorf("config: DB_NAME is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("config: DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("config: DB_PASSWORD is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("config: DB_PORT must be between 1 and 65535, got %d", c.Database.Port)
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("config: SERVER_PORT must be between 1 and 65535, got %d", c.Server.Port)
	}
	switch c.Database.SSLMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		return fmt.Errorf("config: DB_SSLMODE must be one of disable|allow|prefer|require|verify-ca|verify-full, got %q", c.Database.SSLMode)
	}
	return nil
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

func getEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be a valid integer, got %q", key, value)
	}

	return intValue, nil
}

func getEnvDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be a valid duration, got %q", key, value)
	}

	return duration, nil
}

func getEnvBool(name string, defaultVal bool) (bool, error) {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal, nil
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, fmt.Errorf("config: %s must be a valid boolean, got %q", name, valStr)
	}

	return val, nil
}
