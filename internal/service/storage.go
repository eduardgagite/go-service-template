package service

import (
    "context"
	"go-service-template/internal/models"
)

type Storage interface {
	Close() error

    CreateExample(ctx context.Context, example *models.Example) error
    GetExampleByID(ctx context.Context, id int) (*models.Example, error)
    GetAllExamples(ctx context.Context, limit, offset int) ([]models.Example, error)
    UpdateExample(ctx context.Context, example *models.Example) error
    DeleteExample(ctx context.Context, id int) error
}
