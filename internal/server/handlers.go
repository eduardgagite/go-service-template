package server

import (
    "strconv"

    "go-service-template/internal/models"

    "github.com/gofiber/fiber/v2"
)

// healthCheck проверка работоспособности сервиса
// @Summary Health check
// @Description Returns service health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.MessageResponse
// @Router /health [get]
func (s *Server) healthCheck(c *fiber.Ctx) error {
	return c.JSON(models.MessageResponse{
		Message: "Service is healthy",
	})
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
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
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
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
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
		if err.Error() == "example not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
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
		if err.Error() == "example not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
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
		if err.Error() == "example not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(models.MessageResponse{
		Message: "Example deleted successfully",
	})
}
