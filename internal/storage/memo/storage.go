package memo

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
	"github.com/tehrelt/wm-test/internal/storage"
)

type Storage struct {
	m    sync.RWMutex
	data map[uuid.UUID]*models.Task
}

func New() *Storage {
	return &Storage{
		m:    sync.RWMutex{},
		data: make(map[uuid.UUID]*models.Task),
	}
}

func (s *Storage) Save(ctx context.Context, task *models.Task) error {

	s.m.Lock()
	defer s.m.Unlock()
	if _, found := s.data[task.Id]; found {
		task.UpdatedAt = time.Now()
	}

	s.data[task.Id] = task

	return nil
}

func (s *Storage) Find(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	task, found := s.data[id]
	if !found {
		return nil, storage.ErrTaskNotFound
	}

	return task, nil
}

func (s *Storage) List(ctx context.Context) []*models.Task {
	tasks := make([]*models.Task, 0, len(s.data))

	s.m.Lock()
	defer s.m.Unlock()

	if len(s.data) == 0 {
		return tasks
	}

	for _, v := range s.data {
		tasks = append(tasks, v)
	}

	return tasks
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, found := s.data[id]; !found {
		return storage.ErrTaskNotFound
	}

	delete(s.data, id)
	return nil
}
