package config

import (
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DB_PASSWORD", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected DB_PORT=5432, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "service_db" {
		t.Errorf("expected DB_NAME=service_db, got %q", cfg.Database.Name)
	}
	if cfg.Database.User != "postgres" {
		t.Errorf("expected DB_USER=postgres, got %q", cfg.Database.User)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected DB_SSLMODE=disable, got %q", cfg.Database.SSLMode)
	}
	if cfg.Server.Host != "localhost" {
		t.Errorf("expected SERVER_HOST=localhost, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected SERVER_PORT=8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout=10s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 10*time.Second {
		t.Errorf("expected WriteTimeout=10s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.App.DebugMode {
		t.Error("expected DebugMode=false")
	}
	if cfg.App.EnableSwagger {
		t.Error("expected EnableSwagger=false")
	}
	if cfg.Server.BodyLimit != 4*1024*1024 {
		t.Errorf("expected BodyLimit=4MiB, got %d", cfg.Server.BodyLimit)
	}
	if cfg.Server.RateLimit != 100 {
		t.Errorf("expected RateLimit=100, got %d", cfg.Server.RateLimit)
	}
	if cfg.Server.CORSAllowOrigins != "*" {
		t.Errorf("expected CORSAllowOrigins=*, got %q", cfg.Server.CORSAllowOrigins)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_USER", "admin")
	t.Setenv("DB_PASSWORD", "s3cret")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("SERVER_HOST", "0.0.0.0")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("SERVER_READ_TIMEOUT", "30s")
	t.Setenv("SERVER_WRITE_TIMEOUT", "30s")
	t.Setenv("DEBUG_MODE", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("expected db.example.com, got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("expected 5433, got %d", cfg.Database.Port)
	}
	if cfg.Database.Name != "mydb" {
		t.Errorf("expected mydb, got %q", cfg.Database.Name)
	}
	if cfg.Database.User != "admin" {
		t.Errorf("expected admin, got %q", cfg.Database.User)
	}
	if cfg.Database.Password != "s3cret" {
		t.Errorf("expected s3cret, got %q", cfg.Database.Password)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("expected require, got %q", cfg.Database.SSLMode)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected 0.0.0.0, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("expected 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.Server.ReadTimeout)
	}
	if !cfg.App.DebugMode {
		t.Error("expected DebugMode=true")
	}
}

func TestLoad_DatabaseDSN(t *testing.T) {
	t.Setenv("DB_PASSWORD", "pass")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dsn := cfg.DatabaseDSN()
	expected := "host=localhost port=5432 user=postgres password=pass dbname=service_db sslmode=disable"
	if dsn != expected {
		t.Errorf("DSN mismatch:\n  got:  %q\n  want: %q", dsn, expected)
	}
}

func TestLoad_Validation(t *testing.T) {
	t.Run("missing password", func(t *testing.T) {
		// DB_PASSWORD по умолчанию "" — валидация должна это отклонить.
		_, err := Load()
		if err == nil {
			t.Fatal("expected validation error for missing DB_PASSWORD")
		}
	})

	t.Run("empty db name", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("DB_NAME", "")

		// Пустой DB_NAME трактуется getEnv как отсутствующий, поэтому подставляется
		// значение по умолчанию и Load() завершается успешно.
		cfg, err := Load()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Database.Name != "service_db" {
			t.Errorf("expected fallback to service_db, got %q", cfg.Database.Name)
		}
	})

	t.Run("invalid db port", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("DB_PORT", "99999")

		_, err := Load()
		if err == nil {
			t.Fatal("expected validation error for invalid DB_PORT")
		}
	})

	t.Run("invalid server port", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("SERVER_PORT", "0")

		_, err := Load()
		if err == nil {
			t.Fatal("expected validation error for invalid SERVER_PORT")
		}
	})

	t.Run("invalid integer env value", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("DB_MAX_CONNS", "many")

		_, err := Load()
		if err == nil {
			t.Fatal("expected parsing error for invalid DB_MAX_CONNS")
		}
	})

	t.Run("invalid duration env value", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("SERVER_READ_TIMEOUT", "soon")

		_, err := Load()
		if err == nil {
			t.Fatal("expected parsing error for invalid SERVER_READ_TIMEOUT")
		}
	})

	t.Run("invalid boolean env value", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("DEBUG_MODE", "sometimes")

		_, err := Load()
		if err == nil {
			t.Fatal("expected parsing error for invalid DEBUG_MODE")
		}
	})

	t.Run("invalid sslmode", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "pass")
		t.Setenv("DB_SSLMODE", "diasble")

		_, err := Load()
		if err == nil {
			t.Fatal("expected validation error for invalid DB_SSLMODE")
		}
	})
}
