package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
)

func (uc *UseCase) CreateTask(ctx context.Context) (*models.Task, error) {

	logger := uc.logger.With(slog.String("fn", "CreateTask"))

	task := &models.Task{
		Id:        uuid.New(),
		Status:    models.PSWait,
		CreatedAt: time.Now(),
	}

	if err := uc.storage.Save(ctx, task); err != nil {
		logger.Error("failed to save task")
		return nil, err
	}

	uc.processor.Enqueue(task)

	return task, nil
}
