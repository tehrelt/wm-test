package usecase

import (
	"context"
	"time"

	"github.com/tehrelt/wm-test/internal/models"
)

func (uc *UseCase) updateStatus(ctx context.Context, task *models.Task) error {
	if task.Status == models.PSProcessing {
		task.StartProcessingAt = time.Now()
	} else if task.Status == models.PSDone {
		task.EndProcessingAt = time.Now()
	}

	return uc.storage.Save(ctx, task)
}
