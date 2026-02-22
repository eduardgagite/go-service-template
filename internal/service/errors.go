package service

import "errors"

var (
	ErrExampleNotFound     = errors.New("example not found")
	ErrInvalidExampleID    = errors.New("example ID must be positive")
	ErrCreateExampleFailed = errors.New("failed to create example")
	ErrGetExampleFailed    = errors.New("failed to get example")
	ErrGetExamplesFailed   = errors.New("failed to get examples")
	ErrUpdateExampleFailed = errors.New("failed to update example")
	ErrDeleteExampleFailed = errors.New("failed to delete example")
	ErrLimitMustBePositive = errors.New("limit must be positive")
	ErrOffsetMustBeNonNeg  = errors.New("offset must be non-negative")
	ErrRequestCannotBeNil  = errors.New("request cannot be nil")
	ErrNameRequired        = errors.New("name is required")
	ErrNameTooLong         = errors.New("name cannot exceed 255 characters")
	ErrDescriptionTooLong  = errors.New("description cannot exceed 1000 characters")
	ErrValueCannotBeNeg    = errors.New("value cannot be negative")
)
