package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"go-service-template/internal/models"
	storageerrors "go-service-template/internal/storage"
)

type mockStorage struct {
	pingFn          func(ctx context.Context) error
	createExampleFn func(ctx context.Context, example *models.Example) error
	getByIDFn       func(ctx context.Context, id int) (*models.Example, error)
	getAllFn        func(ctx context.Context, limit, offset int) ([]models.Example, error)
	updateFn        func(ctx context.Context, example *models.Example) error
	deleteFn        func(ctx context.Context, id int) error
}

func (m *mockStorage) Ping(ctx context.Context) error {
	if m.pingFn == nil {
		return nil
	}
	return m.pingFn(ctx)
}

func (m *mockStorage) Close() error { return nil }

func (m *mockStorage) CreateExample(ctx context.Context, example *models.Example) error {
	if m.createExampleFn == nil {
		return nil
	}
	return m.createExampleFn(ctx, example)
}

func (m *mockStorage) GetExampleByID(ctx context.Context, id int) (*models.Example, error) {
	if m.getByIDFn == nil {
		return nil, nil
	}
	return m.getByIDFn(ctx, id)
}

func (m *mockStorage) GetAllExamples(ctx context.Context, limit, offset int) ([]models.Example, error) {
	if m.getAllFn == nil {
		return nil, nil
	}
	return m.getAllFn(ctx, limit, offset)
}

func (m *mockStorage) UpdateExample(ctx context.Context, example *models.Example) error {
	if m.updateFn == nil {
		return nil
	}
	return m.updateFn(ctx, example)
}

func (m *mockStorage) DeleteExample(ctx context.Context, id int) error {
	if m.deleteFn == nil {
		return nil
	}
	return m.deleteFn(ctx, id)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestCreateExample_ValidationAndTrimming(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		svc := NewService(&mockStorage{}, testLogger())

		got, err := svc.CreateExample(context.Background(), nil)
		if !errors.Is(err, ErrRequestCannotBeNil) {
			t.Fatalf("expected ErrRequestCannotBeNil, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil example, got: %#v", got)
		}
	})

	t.Run("trims fields and sets timestamps", func(t *testing.T) {
		st := &mockStorage{
			createExampleFn: func(_ context.Context, example *models.Example) error {
				example.ID = 42
				return nil
			},
		}
		svc := NewService(st, testLogger())

		req := &models.ExampleRequest{
			Name:        "  name  ",
			Description: "  desc  ",
			Value:       12.5,
			IsActive:    true,
		}

		got, err := svc.CreateExample(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 42 {
			t.Fatalf("expected ID 42, got: %d", got.ID)
		}
		if got.Name != "name" || got.Description != "desc" {
			t.Fatalf("expected trimmed fields, got name=%q desc=%q", got.Name, got.Description)
		}
		if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
			t.Fatal("expected non-zero timestamps")
		}
	})
}

func TestGetExampleByID_ErrorMapping(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		svc := NewService(&mockStorage{}, testLogger())

		got, err := svc.GetExampleByID(context.Background(), 0)
		if !errors.Is(err, ErrInvalidExampleID) {
			t.Fatalf("expected ErrInvalidExampleID, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil example, got: %#v", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		st := &mockStorage{
			getByIDFn: func(_ context.Context, _ int) (*models.Example, error) {
				return nil, storageerrors.ErrNotFound
			},
		}
		svc := NewService(st, testLogger())

		got, err := svc.GetExampleByID(context.Background(), 1)
		if !errors.Is(err, ErrExampleNotFound) {
			t.Fatalf("expected ErrExampleNotFound, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil example, got: %#v", got)
		}
	})

	t.Run("unexpected storage error", func(t *testing.T) {
		st := &mockStorage{
			getByIDFn: func(_ context.Context, _ int) (*models.Example, error) {
				return nil, errors.New("db down")
			},
		}
		svc := NewService(st, testLogger())

		got, err := svc.GetExampleByID(context.Background(), 1)
		if !errors.Is(err, ErrGetExampleFailed) {
			t.Fatalf("expected ErrGetExampleFailed, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil example, got: %#v", got)
		}
	})
}

func TestGetAllExamples_ValidationAndLimitCap(t *testing.T) {
	t.Run("limit and offset validation", func(t *testing.T) {
		svc := NewService(&mockStorage{}, testLogger())

		if _, err := svc.GetAllExamples(context.Background(), 0, 0); !errors.Is(err, ErrLimitMustBePositive) {
			t.Fatalf("expected ErrLimitMustBePositive, got: %v", err)
		}
		if _, err := svc.GetAllExamples(context.Background(), 1, -1); !errors.Is(err, ErrOffsetMustBeNonNeg) {
			t.Fatalf("expected ErrOffsetMustBeNonNeg, got: %v", err)
		}
	})

	t.Run("caps limit to 100", func(t *testing.T) {
		calledLimit := -1
		calledOffset := -1
		st := &mockStorage{
			getAllFn: func(_ context.Context, limit, offset int) ([]models.Example, error) {
				calledLimit = limit
				calledOffset = offset
				return []models.Example{{ID: 1}}, nil
			},
		}
		svc := NewService(st, testLogger())

		got, err := svc.GetAllExamples(context.Background(), 1000, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if calledLimit != 100 || calledOffset != 5 {
			t.Fatalf("expected storage call with limit=100 offset=5, got limit=%d offset=%d", calledLimit, calledOffset)
		}
		if len(got) != 1 || got[0].ID != 1 {
			t.Fatalf("unexpected data: %#v", got)
		}
	})
}

func TestUpdateExample_ErrorPaths(t *testing.T) {
	t.Run("invalid request data", func(t *testing.T) {
		svc := NewService(&mockStorage{}, testLogger())

		got, err := svc.UpdateExample(context.Background(), 1, &models.ExampleRequest{Name: strings.Repeat("a", 256)})
		if !errors.Is(err, ErrNameTooLong) {
			t.Fatalf("expected ErrNameTooLong, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil example, got: %#v", got)
		}
	})

	t.Run("storage update not found", func(t *testing.T) {
		st := &mockStorage{
			updateFn: func(_ context.Context, _ *models.Example) error {
				return storageerrors.ErrNotFound
			},
		}
		svc := NewService(st, testLogger())

		got, err := svc.UpdateExample(context.Background(), 2, &models.ExampleRequest{Name: "ok"})
		if !errors.Is(err, ErrExampleNotFound) {
			t.Fatalf("expected ErrExampleNotFound, got: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil example, got: %#v", got)
		}
	})

	t.Run("successful update returns loaded record", func(t *testing.T) {
		now := time.Now()
		st := &mockStorage{
			updateFn: func(_ context.Context, _ *models.Example) error {
				return nil
			},
			getByIDFn: func(_ context.Context, id int) (*models.Example, error) {
				return &models.Example{ID: id, Name: "n", UpdatedAt: now}, nil
			},
		}
		svc := NewService(st, testLogger())

		got, err := svc.UpdateExample(context.Background(), 2, &models.ExampleRequest{Name: " n "})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.ID != 2 || got.Name != "n" {
			t.Fatalf("unexpected updated example: %#v", got)
		}
	})
}

func TestDeleteExample_ErrorMapping(t *testing.T) {
	st := &mockStorage{
		deleteFn: func(_ context.Context, _ int) error {
			return storageerrors.ErrNotFound
		},
	}
	svc := NewService(st, testLogger())

	err := svc.DeleteExample(context.Background(), 5)
	if !errors.Is(err, ErrExampleNotFound) {
		t.Fatalf("expected ErrExampleNotFound, got: %v", err)
	}
}
