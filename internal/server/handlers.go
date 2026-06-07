package server

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go-service-template/internal/models"
	"go-service-template/internal/service"

	"github.com/gofiber/fiber/v2"
)

// liveness проверка, что процесс жив (без внешних зависимостей)
// @Summary Liveness probe
// @Description Returns 200 while the process is running. Wire to k8s livenessProbe.
// @Tags health
// @Produce json
// @Success 200 {object} models.MessageResponse
// @Router /livez [get]
func (s *Server) liveness(c *fiber.Ctx) error {
	return c.JSON(models.MessageResponse{Message: "alive"})
}

// readiness проверка готовности обслуживать трафик (доступность БД)
// @Summary Readiness probe
// @Description Returns 200 when dependencies are reachable, 503 otherwise. Wire to k8s readinessProbe.
// @Tags health
// @Produce json
// @Success 200 {object} models.MessageResponse
// @Failure 503 {object} models.ErrorResponse "Service unavailable"
// @Router /readyz [get]
func (s *Server) readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	if err := s.services.Ping(ctx); err != nil {
		s.logger.Error("readiness check failed", "error", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(models.ErrorResponse{
			Error: "database is unavailable",
		})
	}

	return c.JSON(models.MessageResponse{Message: "ready"})
}

// createExample создает новый пример
// @Summary Create example
// @Description Creates a new example record
// @Tags examples
// @Accept json
// @Produce json
// @Param example body models.ExampleRequest true "Example data"
// @Success 201 {object} models.Example
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Router /examples [post]
func (s *Server) createExample(c *fiber.Ctx) error {
	var req models.ExampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
	}

	example, err := s.services.Example.CreateExample(c.UserContext(), &req)
	if err != nil {
		return s.handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(example)
}

// getAllExamples получает список всех примеров
// @Summary Get all examples
// @Description Returns a list of all examples with pagination
// @Tags examples
// @Accept json
// @Produce json
// @Param limit query int false "Number of records" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} models.ExampleResponse
// @Failure 400 {object} models.ErrorResponse "Invalid parameters"
// @Router /examples [get]
func (s *Server) getAllExamples(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid limit parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid offset parameter",
		})
	}

	examples, err := s.services.Example.GetAllExamples(c.UserContext(), limit, offset)
	if err != nil {
		return s.handleServiceError(c, err)
	}

	return c.JSON(models.ExampleResponse{
		Data: examples,
	})
}

// getExample получает пример по ID
// @Summary Get example by ID
// @Description Returns an example by its ID
// @Tags examples
// @Accept json
// @Produce json
// @Param id path int true "Example ID"
// @Success 200 {object} models.Example
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "Example not found"
// @Router /examples/{id} [get]
func (s *Server) getExample(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid example ID",
		})
	}

	example, err := s.services.Example.GetExampleByID(c.UserContext(), id)
	if err != nil {
		return s.handleServiceError(c, err)
	}

	return c.JSON(example)
}

// updateExample обновляет пример
// @Summary Update example
// @Description Updates an existing example
// @Tags examples
// @Accept json
// @Produce json
// @Param id path int true "Example ID"
// @Param example body models.ExampleRequest true "Updated example data"
// @Success 200 {object} models.Example
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Failure 404 {object} models.ErrorResponse "Example not found"
// @Router /examples/{id} [put]
func (s *Server) updateExample(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid example ID",
		})
	}

	var req models.ExampleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
	}

	example, err := s.services.Example.UpdateExample(c.UserContext(), id, &req)
	if err != nil {
		return s.handleServiceError(c, err)
	}

	return c.JSON(example)
}

// deleteExample удаляет пример
// @Summary Delete example
// @Description Deletes an example by ID
// @Tags examples
// @Accept json
// @Produce json
// @Param id path int true "Example ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse "Invalid ID"
// @Failure 404 {object} models.ErrorResponse "Example not found"
// @Router /examples/{id} [delete]
func (s *Server) deleteExample(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Invalid example ID",
		})
	}

	if err := s.services.Example.DeleteExample(c.UserContext(), id); err != nil {
		return s.handleServiceError(c, err)
	}

	return c.JSON(models.MessageResponse{
		Message: "Example deleted successfully",
	})
}

func (s *Server) handleServiceError(c *fiber.Ctx, err error) error {
	status := mapServiceErrorToHTTPStatus(err)
	if status == fiber.StatusInternalServerError {
		s.logger.Error("Unhandled service error", "error", err)
		return c.Status(status).JSON(models.ErrorResponse{
			Error: "internal server error",
		})
	}

	return c.Status(status).JSON(models.ErrorResponse{
		Error: err.Error(),
	})
}

func mapServiceErrorToHTTPStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrExampleNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, service.ErrInvalidExampleID),
		errors.Is(err, service.ErrLimitMustBePositive),
		errors.Is(err, service.ErrOffsetMustBeNonNeg),
		errors.Is(err, service.ErrRequestCannotBeNil),
		errors.Is(err, service.ErrNameRequired),
		errors.Is(err, service.ErrNameTooLong),
		errors.Is(err, service.ErrDescriptionTooLong),
		errors.Is(err, service.ErrValueCannotBeNeg):
		return fiber.StatusBadRequest
	default:
		return fiber.StatusInternalServerError
	}
}
