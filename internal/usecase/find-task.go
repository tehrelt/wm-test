package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
	"github.com/tehrelt/wm-test/internal/storage"
)

func (uc *UseCase) FindTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	logger := uc.logger.With(slog.String("fn", "FindTask"))

	task, err := uc.storage.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			logger.Warn("task not found", slog.String("id", id.String()))
			return nil, err
		}

		logger.Error("failed to save task")
		return nil, err
	}

	return task, nil
}
