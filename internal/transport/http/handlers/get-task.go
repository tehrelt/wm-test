package handlers

import (
	"errors"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/storage"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers/dto"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func GetTask(uc *usecase.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) {

		paramId := c.Params("id")

		id, err := uuid.Parse(paramId)
		if err != nil {
			c.
				Status(fiber.StatusBadRequest).
				JSON(fiber.Map{
					"error":   "invalid id",
					"details": err.Error(),
				})
			return
		}

		task, err := uc.FindTask(c.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				c.Status(404).JSON(fiber.Map{
					"error": "task not found",
					"id":    id.String(),
				})
				return
			}
		}

		c.JSON(dto.TaskFrom(task))
	}
}
