package models

import (
	"time"

	"github.com/google/uuid"
)

type ProcessStatus int

func (ps ProcessStatus) String() string {
	labels := []string{"nil", "wait", "processing", "done", "error"}

	if ps < 0 || int(ps) > len(labels) {
		return "undefined"
	}

	return labels[ps]
}

const (
	PSNil ProcessStatus = iota
	PSWait
	PSProcessing
	PSDone
	PSError
)

type Task struct {
	Id        uuid.UUID
	Status    ProcessStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (t Task) Elapsed() time.Duration {
	if t.UpdatedAt != nil {
		return t.UpdatedAt.Sub(t.CreatedAt)
	}

	return time.Since(t.CreatedAt)
}
