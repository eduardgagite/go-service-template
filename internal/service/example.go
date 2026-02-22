package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"go-service-template/internal/models"
	storageerrors "go-service-template/internal/storage"
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

func (s *service) CreateExample(ctx context.Context, req *models.ExampleRequest) (*models.Example, error) {
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

	if err := s.storage.CreateExample(ctx, example); err != nil {
		s.logger.Error("Failed to create example", slog.String("error", err.Error()))
		return nil, ErrCreateExampleFailed
	}

	s.logger.Info("Example created successfully", slog.Int("id", example.ID))
	return example, nil
}

func (s *service) GetExampleByID(ctx context.Context, id int) (*models.Example, error) {
	if id <= 0 {
		return nil, ErrInvalidExampleID
	}

	example, err := s.storage.GetExampleByID(ctx, id)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			return nil, ErrExampleNotFound
		}
		s.logger.Error("Failed to get example", slog.Int("id", id), slog.String("error", err.Error()))
		return nil, ErrGetExampleFailed
	}

	return example, nil
}

func (s *service) GetAllExamples(ctx context.Context, limit, offset int) ([]models.Example, error) {
	if limit <= 0 {
		return nil, ErrLimitMustBePositive
	}
	if offset < 0 {
		return nil, ErrOffsetMustBeNonNeg
	}
	if limit > 100 {
		limit = 100
	}

	examples, err := s.storage.GetAllExamples(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get examples", slog.String("error", err.Error()))
		return nil, ErrGetExamplesFailed
	}

	return examples, nil
}

func (s *service) UpdateExample(ctx context.Context, id int, req *models.ExampleRequest) (*models.Example, error) {
	if id <= 0 {
		return nil, ErrInvalidExampleID
	}

	if err := s.validateExampleRequest(req); err != nil {
		return nil, err
	}

	example := &models.Example{
		ID:          id,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Value:       req.Value,
		IsActive:    req.IsActive,
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.UpdateExample(ctx, example); err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			return nil, ErrExampleNotFound
		}
		s.logger.Error("Failed to update example", slog.Int("id", id), slog.String("error", err.Error()))
		return nil, ErrUpdateExampleFailed
	}

	s.logger.Info("Example updated successfully", slog.Int("id", id))
	updatedExample, err := s.storage.GetExampleByID(ctx, id)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			return nil, ErrExampleNotFound
		}
		s.logger.Error("Failed to load updated example", slog.Int("id", id), slog.String("error", err.Error()))
		return nil, ErrGetExampleFailed
	}
	return updatedExample, nil
}

func (s *service) DeleteExample(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidExampleID
	}

	if err := s.storage.DeleteExample(ctx, id); err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			return ErrExampleNotFound
		}
		s.logger.Error("Failed to delete example", slog.Int("id", id), slog.String("error", err.Error()))
		return ErrDeleteExampleFailed
	}

	s.logger.Info("Example deleted successfully", slog.Int("id", id))
	return nil
}

func (s *service) validateExampleRequest(req *models.ExampleRequest) error {
	if req == nil {
		return ErrRequestCannotBeNil
	}

	if strings.TrimSpace(req.Name) == "" {
		return ErrNameRequired
	}

	if len(req.Name) > 255 {
		return ErrNameTooLong
	}

	if len(req.Description) > 1000 {
		return ErrDescriptionTooLong
	}

	if req.Value < 0 {
		return ErrValueCannotBeNeg
	}

	return nil
}
