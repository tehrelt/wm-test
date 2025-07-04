package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tehrelt/wm-test/internal/storage"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func DeleteTask(uc *usecase.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		paramId := c.Param("id")

		id, err := uuid.Parse(paramId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = uc.DeleteTask(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "task not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
