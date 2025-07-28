package models

import (
	"time"
)

type Example struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Value       float64   `json:"value"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExampleRequest struct {
	Name        string  `json:"name" example:"Example Name"`
	Description string  `json:"description" example:"Example description"`
	Value       float64 `json:"value" example:"100.50"`
	IsActive    bool    `json:"is_active" example:"true"`
}

type ExampleResponse struct {
	Data []Example `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input data"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
