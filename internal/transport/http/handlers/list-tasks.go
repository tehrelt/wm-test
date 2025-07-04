package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers/dto"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func ListTasks(uc *usecase.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {

		tasks, _ := uc.ListTasks(c.Request().Context())

		response := struct {
			Tasks []*dto.Task `json:"tasks"`
			Total uint        `json:"total"`
		}{
			Tasks: make([]*dto.Task, len(tasks)),
			Total: uint(len(tasks)),
		}

		for i := range tasks {
			response.Tasks[i] = dto.TaskFrom(tasks[i])
		}

		return c.JSON(http.StatusOK, response)
	}
}
