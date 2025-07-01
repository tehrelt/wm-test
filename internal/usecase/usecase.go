package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
)

type TaskStorage interface {
	Save(context.Context, *models.Task) error
	Find(context.Context, uuid.UUID) (*models.Task, error)
	List(context.Context) []*models.Task
	Delete(context.Context, uuid.UUID) error
}

type UseCase struct {
	storage TaskStorage
	logger  *slog.Logger
}

func New(storage TaskStorage) *UseCase {
	return &UseCase{
		storage: storage,
		logger:  slog.With(slog.String("comp", "usecase.UseCase")),
	}
}
