package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/storage"
)

func (uc *UseCase) DeleteTask(ctx context.Context, id uuid.UUID) error {
	logger := uc.logger.With(slog.String("fn", "DeleteTask"))

	task, err := uc.storage.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			logger.Warn("task not found", slog.String("id", id.String()))
			return err
		}

		logger.Error("failed to find task", slog.String("id", id.String()), slog.String("err", err.Error()))
		return err
	}

	if task.Status.IsCancelable() {
		logger.Info("canceling task", slog.String("id", id.String()))
		uc.processor.Cancel(id)
	}

	if err := uc.storage.Delete(ctx, id); err != nil {
		logger.Error("failed to delete task", slog.String("id", id.String()), slog.String("err", err.Error()))
		return err
	}

	return nil
}
