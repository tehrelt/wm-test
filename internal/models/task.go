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

func (ps ProcessStatus) IsCancelable() bool {
	return ps == PSWait || ps == PSProcessing
}

const (
	PSNil ProcessStatus = iota
	PSWait
	PSProcessing
	PSDone
	PSError
)

type Task struct {
	Id                uuid.UUID
	Status            ProcessStatus
	StartProcessingAt time.Time
	EndProcessingAt   time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (t Task) Elapsed() time.Duration {
	if t.Status == PSDone {
		return t.EndProcessingAt.Sub(t.StartProcessingAt)
	}

	if t.Status == PSProcessing {
		return time.Since(t.StartProcessingAt)
	}

	return 0
}

func (t Task) IsCancelable() bool {
	return t.Status == PSWait || t.Status == PSProcessing
}
