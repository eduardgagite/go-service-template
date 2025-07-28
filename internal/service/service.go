package service

import (
	"log/slog"

	"go-service-template/internal/models"
)

type Service interface {
	CreateExample(req *models.ExampleRequest) (*models.Example, error)
	GetExampleByID(id int) (*models.Example, error)
	GetAllExamples(limit, offset int) ([]models.Example, error)
	UpdateExample(id int, req *models.ExampleRequest) (*models.Example, error)
	DeleteExample(id int) error
}

type Services struct {
	Example Service
}

func NewServices(storage Storage, logger *slog.Logger) *Services {
	return &Services{
		Example: NewService(storage, logger),
	}
}
