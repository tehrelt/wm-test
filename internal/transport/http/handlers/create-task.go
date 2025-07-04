package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers/dto"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func CreateTask(uc *usecase.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {

		task, err := uc.CreateTask(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, dto.TaskFrom(task))
	}
}
