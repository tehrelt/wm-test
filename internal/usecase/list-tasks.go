package usecase

import (
	"context"

	"github.com/tehrelt/wm-test/internal/models"
)

func (uc *UseCase) ListTasks(ctx context.Context) ([]*models.Task, error) {

	tasks := uc.storage.List(ctx)

	return tasks, nil
}
