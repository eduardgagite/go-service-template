package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-service-template/internal/config"
	"go-service-template/internal/models"
	"go-service-template/internal/service"
)

type mockExampleService struct {
	createFn  func(ctx context.Context, req *models.ExampleRequest) (*models.Example, error)
	getByIDFn func(ctx context.Context, id int) (*models.Example, error)
	getAllFn  func(ctx context.Context, limit, offset int) ([]models.Example, error)
	updateFn  func(ctx context.Context, id int, req *models.ExampleRequest) (*models.Example, error)
	deleteFn  func(ctx context.Context, id int) error
}

func (m *mockExampleService) CreateExample(ctx context.Context, req *models.ExampleRequest) (*models.Example, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m *mockExampleService) GetExampleByID(ctx context.Context, id int) (*models.Example, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockExampleService) GetAllExamples(ctx context.Context, limit, offset int) ([]models.Example, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockExampleService) UpdateExample(ctx context.Context, id int, req *models.ExampleRequest) (*models.Example, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockExampleService) DeleteExample(ctx context.Context, id int) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func newTestServer(mock *mockExampleService, pingFn func(ctx context.Context) error) *Server {
	services := &service.Services{
		Example:  mock,
		PingFunc: pingFn,
	}
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s := New(services, logger, cfg)
	s.setupRoutes()
	return s
}

func doRequest(s *Server, method, path string, body any) *http.Response {
	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := s.app.Test(req, -1)
	return resp
}

func decodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return v
}

func TestLiveness(t *testing.T) {
	// Liveness не должен зависеть от БД: падающий ping всё равно даёт 200.
	s := newTestServer(&mockExampleService{}, func(ctx context.Context) error {
		return errors.New("db down")
	})

	resp := doRequest(s, http.MethodGet, "/livez", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body := decodeJSON[models.MessageResponse](t, resp)
	if body.Message != "alive" {
		t.Fatalf("unexpected message: %q", body.Message)
	}
}

func TestReadiness(t *testing.T) {
	t.Run("ready", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, func(ctx context.Context) error { return nil })

		resp := doRequest(s, http.MethodGet, "/readyz", nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body := decodeJSON[models.MessageResponse](t, resp)
		if body.Message != "ready" {
			t.Fatalf("unexpected message: %q", body.Message)
		}
	})

	t.Run("database unavailable", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, func(ctx context.Context) error {
			return errors.New("connection refused")
		})

		resp := doRequest(s, http.MethodGet, "/readyz", nil)
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d", resp.StatusCode)
		}
		body := decodeJSON[models.ErrorResponse](t, resp)
		if body.Error != "database is unavailable" {
			t.Fatalf("unexpected error: %q", body.Error)
		}
	})

	t.Run("alias /health maps to readiness", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, func(ctx context.Context) error { return nil })

		resp := doRequest(s, http.MethodGet, "/health", nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
}

func TestCreateExample(t *testing.T) {
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock := &mockExampleService{
			createFn: func(_ context.Context, req *models.ExampleRequest) (*models.Example, error) {
				return &models.Example{ID: 1, Name: req.Name, CreatedAt: now, UpdatedAt: now}, nil
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodPost, "/api/v1/examples", models.ExampleRequest{
			Name: "test", Value: 10,
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}
		body := decodeJSON[models.Example](t, resp)
		if body.ID != 1 || body.Name != "test" {
			t.Fatalf("unexpected body: %+v", body)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/examples", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := s.app.Test(req, -1)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("service validation error", func(t *testing.T) {
		mock := &mockExampleService{
			createFn: func(_ context.Context, _ *models.ExampleRequest) (*models.Example, error) {
				return nil, service.ErrNameRequired
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodPost, "/api/v1/examples", models.ExampleRequest{})
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
}

func TestGetAllExamples(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockExampleService{
			getAllFn: func(_ context.Context, limit, offset int) ([]models.Example, error) {
				return []models.Example{{ID: 1}, {ID: 2}}, nil
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodGet, "/api/v1/examples?limit=10&offset=0", nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body := decodeJSON[models.ExampleResponse](t, resp)
		if len(body.Data) != 2 {
			t.Fatalf("expected 2 items, got %d", len(body.Data))
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		resp := doRequest(s, http.MethodGet, "/api/v1/examples?limit=abc", nil)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid offset", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		resp := doRequest(s, http.MethodGet, "/api/v1/examples?offset=abc", nil)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
}

func TestGetExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockExampleService{
			getByIDFn: func(_ context.Context, id int) (*models.Example, error) {
				return &models.Example{ID: id, Name: "found"}, nil
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodGet, "/api/v1/examples/5", nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body := decodeJSON[models.Example](t, resp)
		if body.ID != 5 {
			t.Fatalf("expected ID 5, got %d", body.ID)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		resp := doRequest(s, http.MethodGet, "/api/v1/examples/abc", nil)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockExampleService{
			getByIDFn: func(_ context.Context, _ int) (*models.Example, error) {
				return nil, service.ErrExampleNotFound
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodGet, "/api/v1/examples/99", nil)
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockExampleService{
			updateFn: func(_ context.Context, id int, req *models.ExampleRequest) (*models.Example, error) {
				return &models.Example{ID: id, Name: req.Name}, nil
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodPut, "/api/v1/examples/3", models.ExampleRequest{Name: "upd"})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body := decodeJSON[models.Example](t, resp)
		if body.ID != 3 || body.Name != "upd" {
			t.Fatalf("unexpected body: %+v", body)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		resp := doRequest(s, http.MethodPut, "/api/v1/examples/abc", models.ExampleRequest{Name: "x"})
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/examples/1", bytes.NewReader([]byte("bad")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := s.app.Test(req, -1)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockExampleService{
			updateFn: func(_ context.Context, _ int, _ *models.ExampleRequest) (*models.Example, error) {
				return nil, service.ErrExampleNotFound
			},
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodPut, "/api/v1/examples/99", models.ExampleRequest{Name: "x"})
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockExampleService{
			deleteFn: func(_ context.Context, _ int) error { return nil },
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodDelete, "/api/v1/examples/1", nil)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		s := newTestServer(&mockExampleService{}, nil)

		resp := doRequest(s, http.MethodDelete, "/api/v1/examples/abc", nil)
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockExampleService{
			deleteFn: func(_ context.Context, _ int) error { return service.ErrExampleNotFound },
		}
		s := newTestServer(mock, nil)

		resp := doRequest(s, http.MethodDelete, "/api/v1/examples/99", nil)
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
}

func TestMapServiceErrorToHTTPStatus(t *testing.T) {
	tests := []struct {
		err    error
		status int
	}{
		{service.ErrExampleNotFound, 404},
		{service.ErrInvalidExampleID, 400},
		{service.ErrLimitMustBePositive, 400},
		{service.ErrOffsetMustBeNonNeg, 400},
		{service.ErrRequestCannotBeNil, 400},
		{service.ErrNameRequired, 400},
		{service.ErrNameTooLong, 400},
		{service.ErrDescriptionTooLong, 400},
		{service.ErrValueCannotBeNeg, 400},
		{service.ErrCreateExampleFailed, 500},
		{errors.New("unknown"), 500},
	}

	for _, tt := range tests {
		got := mapServiceErrorToHTTPStatus(tt.err)
		if got != tt.status {
			t.Errorf("mapServiceErrorToHTTPStatus(%v) = %d, want %d", tt.err, got, tt.status)
		}
	}
}
