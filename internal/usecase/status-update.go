package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/tehrelt/wm-test/internal/models"
)

func (uc *UseCase) updateStatus(ctx context.Context, task *models.Task) error {

	if task.Status == models.PSProcessing {
		task.StartProcessingAt = time.Now()
	} else if task.Status == models.PSDone {
		task.EndProcessingAt = time.Now()
	}

	uc.logger.Info("updating task status", slog.String("task_id", task.Id.String()), slog.String("status", task.Status.String()))

	return uc.storage.Save(ctx, task)
}
