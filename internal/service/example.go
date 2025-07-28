package service

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"go-service-template/internal/models"
)

type service struct {
	storage Storage
	logger  *slog.Logger
}

func NewService(storage Storage, logger *slog.Logger) Service {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

func (s *service) CreateExample(req *models.ExampleRequest) (*models.Example, error) {
	if err := s.validateExampleRequest(req); err != nil {
		return nil, err
	}

	example := &models.Example{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Value:       req.Value,
		IsActive:    req.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.CreateExample(example); err != nil {
		s.logger.Error("Failed to create example", slog.String("error", err.Error()))
		return nil, errors.New("failed to create example")
	}

	s.logger.Info("Example created successfully", slog.Int("id", example.ID))
	return example, nil
}

func (s *service) GetExampleByID(id int) (*models.Example, error) {
	if id <= 0 {
		return nil, errors.New("example ID must be positive")
	}

	example, err := s.storage.GetExampleByID(id)
	if err != nil {
		s.logger.Error("Failed to get example", slog.Int("id", id), slog.String("error", err.Error()))
		return nil, errors.New("failed to get example")
	}

	if example == nil {
		return nil, errors.New("example not found")
	}

	return example, nil
}

func (s *service) GetAllExamples(limit, offset int) ([]models.Example, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	if offset < 0 {
		return nil, errors.New("offset must be non-negative")
	}
	if limit > 100 {
		limit = 100
	}

	examples, err := s.storage.GetAllExamples(limit, offset)
	if err != nil {
		s.logger.Error("Failed to get examples", slog.String("error", err.Error()))
		return nil, errors.New("failed to get examples")
	}

	return examples, nil
}

func (s *service) UpdateExample(id int, req *models.ExampleRequest) (*models.Example, error) {
	if id <= 0 {
		return nil, errors.New("example ID must be positive")
	}

	if err := s.validateExampleRequest(req); err != nil {
		return nil, err
	}

	existingExample, err := s.storage.GetExampleByID(id)
	if err != nil {
		s.logger.Error("Failed to get example for update", slog.Int("id", id), slog.String("error", err.Error()))
		return nil, errors.New("failed to get example")
	}
	if existingExample == nil {
		return nil, errors.New("example not found")
	}

	example := &models.Example{
		ID:          id,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Value:       req.Value,
		IsActive:    req.IsActive,
		CreatedAt:   existingExample.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.UpdateExample(example); err != nil {
		s.logger.Error("Failed to update example", slog.Int("id", id), slog.String("error", err.Error()))
		return nil, errors.New("failed to update example")
	}

	s.logger.Info("Example updated successfully", slog.Int("id", id))
	return example, nil
}

func (s *service) DeleteExample(id int) error {
	if id <= 0 {
		return errors.New("example ID must be positive")
	}

	existingExample, err := s.storage.GetExampleByID(id)
	if err != nil {
		s.logger.Error("Failed to get example for deletion", slog.Int("id", id), slog.String("error", err.Error()))
		return errors.New("failed to get example")
	}
	if existingExample == nil {
		return errors.New("example not found")
	}

	if err := s.storage.DeleteExample(id); err != nil {
		s.logger.Error("Failed to delete example", slog.Int("id", id), slog.String("error", err.Error()))
		return errors.New("failed to delete example")
	}

	s.logger.Info("Example deleted successfully", slog.Int("id", id))
	return nil
}

func (s *service) validateExampleRequest(req *models.ExampleRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	if len(req.Name) > 255 {
		return errors.New("name cannot exceed 255 characters")
	}

	if len(req.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	if req.Value < 0 {
		return errors.New("value cannot be negative")
	}

	return nil
}
