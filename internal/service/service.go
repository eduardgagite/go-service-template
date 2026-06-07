package service

import (
	"context"
	"log/slog"

	"go-service-template/internal/models"
)

type Service interface {
	CreateExample(ctx context.Context, req *models.ExampleRequest) (*models.Example, error)
	GetExampleByID(ctx context.Context, id int) (*models.Example, error)
	GetAllExamples(ctx context.Context, limit, offset int) ([]models.Example, error)
	UpdateExample(ctx context.Context, id int, req *models.ExampleRequest) (*models.Example, error)
	DeleteExample(ctx context.Context, id int) error
}

type Services struct {
	Example  Service
	PingFunc func(ctx context.Context) error
}

func NewServices(storage Storage, logger *slog.Logger) *Services {
	return &Services{
		Example:  NewService(storage, logger),
		PingFunc: storage.Ping,
	}
}

func (s *Services) Ping(ctx context.Context) error {
	if s.PingFunc == nil {
		return nil
	}
	return s.PingFunc(ctx)
}
