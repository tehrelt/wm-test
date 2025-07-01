package handlers

import (
	"github.com/gofiber/fiber"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers/dto"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func CreateTask(uc *usecase.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) {

		task, err := uc.CreateTask(c.Context())
		if err != nil {
			c.Status(500).JSON(fiber.Map{
				"error":   "unexpected error",
				"details": err.Error(),
			})

			return
		}

		c.JSON(dto.TaskFrom(task))
	}
}
