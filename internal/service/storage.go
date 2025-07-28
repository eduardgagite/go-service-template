package service

import (
	"go-service-template/internal/models"
)

type Storage interface {
	Close() error

	CreateExample(example *models.Example) error
	GetExampleByID(id int) (*models.Example, error)
	GetAllExamples(limit, offset int) ([]models.Example, error)
	UpdateExample(example *models.Example) error
	DeleteExample(id int) error
}
