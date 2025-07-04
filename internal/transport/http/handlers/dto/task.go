package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
)

type Task struct {
	Id        uuid.UUID     `json:"id"`
	Status    string        `json:"status"`
	Elapsed   time.Duration `json:"elapsed"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

func TaskFrom(m *models.Task) *Task {

	t := &Task{}

	t.Id = m.Id
	t.Status = m.Status.String()
	t.CreatedAt = m.CreatedAt
	t.Elapsed = m.Elapsed()
	t.UpdatedAt = m.UpdatedAt

	return t
}
