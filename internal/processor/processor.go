package processor

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/tehrelt/wm-test/internal/models"
)

type StatusUpdateFn func(context.Context, *models.Task) error

type processingTask struct {
	*models.Task
	update StatusUpdateFn
	ctx    context.Context
}

func (pt *processingTask) SetStatus(ctx context.Context, status models.ProcessStatus) error {
	pt.Status = status

	if err := pt.update(ctx, pt.Task); err != nil {
		return err
	}

	return nil
}

type processingTaskFactory struct {
	update StatusUpdateFn
}

func (ptf *processingTaskFactory) New(task *models.Task) *processingTask {
	return &processingTask{
		Task:   task,
		update: ptf.update,
	}
}

func (ptf *processingTaskFactory) NewWithContext(task *models.Task, ctx context.Context) *processingTask {
	return &processingTask{
		Task:   task,
		update: ptf.update,
		ctx:    ctx,
	}
}

type cancelMap struct {
	m    sync.Mutex
	data map[uuid.UUID]context.CancelFunc
}

func newCancelMap() *cancelMap {
	return &cancelMap{data: make(map[uuid.UUID]context.CancelFunc)}
}

func (cm *cancelMap) set(id uuid.UUID, cancel context.CancelFunc) {
	cm.m.Lock()
	defer cm.m.Unlock()
	cm.data[id] = cancel
}

func (cm *cancelMap) delete(id uuid.UUID) {
	cm.m.Lock()
	defer cm.m.Unlock()
	delete(cm.data, id)
}

func (cm *cancelMap) get(id uuid.UUID) (context.CancelFunc, bool) {
	cm.m.Lock()
	defer cm.m.Unlock()
	cancel, ok := cm.data[id]
	return cancel, ok
}

type TaskProcessor struct {
	queue   chan *processingTask
	factory *processingTaskFactory
	logger  *slog.Logger

	cancels *cancelMap
}

func NewTaskProcessor(update StatusUpdateFn, queueSize int) *TaskProcessor {
	tp := &TaskProcessor{
		queue:   make(chan *processingTask, queueSize),
		factory: &processingTaskFactory{update: update},
		logger:  slog.With(slog.String("comp", "processor.TaskProcessor")),
		cancels: newCancelMap(),
	}

	return tp
}

func (tp *TaskProcessor) Start(ctx context.Context, workerCount int) error {
	if workerCount <= 0 {
		return errors.New("workerCount must be greater than 0")
	}

	for i := 0; i < workerCount; i++ {
		w := newWorker(i, tp.queue)
		go w.run(ctx)
	}

	return nil
}

func (tp *TaskProcessor) Enqueue(task *models.Task) {
	ctx, cancel := context.WithCancel(context.Background())
	tp.cancels.set(task.Id, cancel)
	tp.queue <- tp.factory.NewWithContext(task, ctx)
}

func (tp *TaskProcessor) Cancel(taskID uuid.UUID) bool {
	cancel, found := tp.cancels.get(taskID)
	if found {
		cancel()
		tp.cancels.delete(taskID)
		return true
	}

	return false
}
