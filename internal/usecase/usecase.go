package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/config"
	"github.com/tehrelt/wm-test/internal/models"
	"github.com/tehrelt/wm-test/internal/processor"
)

type TaskStorage interface {
	Save(context.Context, *models.Task) error
	Find(context.Context, uuid.UUID) (*models.Task, error)
	List(context.Context) []*models.Task
	Delete(context.Context, uuid.UUID) error
}

type UseCase struct {
	cfg       *config.Config
	storage   TaskStorage
	logger    *slog.Logger
	processor *processor.TaskProcessor
}

func New(cfg *config.Config, storage TaskStorage) *UseCase {
	uc := &UseCase{
		cfg:       cfg,
		storage:   storage,
		processor: processor.NewTaskProcessor(100),
		logger:    slog.With(slog.String("comp", "usecase.UseCase")),
	}

	uc.processor.Start(context.Background(), uc.cfg.WorkerCount)
	uc.processor.Subscribe(uc.updateStatus)

	return uc
}
