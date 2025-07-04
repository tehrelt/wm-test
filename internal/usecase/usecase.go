package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
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
	storage   TaskStorage
	logger    *slog.Logger
	processor *processor.TaskProcessor
}

func New(storage TaskStorage) *UseCase {
	uc := &UseCase{
		storage: storage,
		logger:  slog.With(slog.String("comp", "usecase.UseCase")),
	}

	uc.setup(context.Background())

	return uc
}

func (uc *UseCase) setup(ctx context.Context) {
	uc.processor = processor.NewTaskProcessor(uc.updateStatus, 100)
	uc.processor.Start(ctx, 4)
}
