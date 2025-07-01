package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers/dto"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func ListTasks(uc *usecase.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) {

		tasks, _ := uc.ListTasks(c.Context())
		slog.Debug("list of tasks", slog.Any("list", tasks))

		response := struct {
			Tasks []*dto.Task `json:"tasks"`
			Total uint        `json:"total"`
		}{
			Tasks: make([]*dto.Task, len(tasks)),
			Total: uint(len(tasks)),
		}
		slog.Info("response struct", slog.Any("response", response))

		for i := range tasks {
			response.Tasks[i] = dto.TaskFrom(tasks[i])
			slog.Info("task mapped", slog.Any("task", response.Tasks[i]))
		}

		slog.Info("response struct after mapping", slog.Any("response", response))

		c.JSON(response)
	}
}
