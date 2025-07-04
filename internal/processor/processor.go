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
	ch  chan<- *models.Task
	ctx context.Context
}

func (pt *processingTask) SetStatus(ctx context.Context, status models.ProcessStatus) error {
	pt.Status = status

	pt.ch <- pt.Task

	return nil
}

type processingTaskFactory struct {
	pending chan *models.Task
}

func newProcessingTaskFactory(queueSize int) *processingTaskFactory {
	return &processingTaskFactory{pending: make(chan *models.Task, queueSize)}
}

func (ptf *processingTaskFactory) NewWithContext(ctx context.Context, task *models.Task) *processingTask {
	return &processingTask{
		ctx:  ctx,
		Task: task,
		ch:   ptf.pending,
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

type onStatusUpdateFn func(context.Context, *models.Task) error

type TaskProcessor struct {
	queue   chan *processingTask
	factory *processingTaskFactory
	logger  *slog.Logger

	subs []onStatusUpdateFn

	cancels *cancelMap
}

func NewTaskProcessor(queueSize int) *TaskProcessor {
	tp := &TaskProcessor{
		queue:   make(chan *processingTask, queueSize),
		logger:  slog.With(slog.String("comp", "processor.TaskProcessor")),
		factory: newProcessingTaskFactory(queueSize),
		cancels: newCancelMap(),
		subs:    make([]onStatusUpdateFn, 0),
	}

	return tp
}

func (tp *TaskProcessor) Start(ctx context.Context, workerCount int) error {
	if workerCount <= 0 {
		return errors.New("workerCount must be greater than 0")
	}

	tp.logger.Info("starting task processor", slog.Int("worker_count", workerCount))
	for i := 0; i < workerCount; i++ {
		w := newWorker(i, tp.queue)
		go w.run(ctx)
	}

	go func() {
		for task := range tp.factory.pending {
			for _, fn := range tp.subs {
				if err := fn(ctx, task); err != nil {
					tp.logger.Error("failed to update task status", "error", err)
				}
			}
		}
	}()

	return nil
}

func (tp *TaskProcessor) Enqueue(task *models.Task) {
	ctx, cancel := context.WithCancel(context.Background())
	tp.cancels.set(task.Id, cancel)
	tp.queue <- tp.factory.NewWithContext(ctx, task)
}

func (tp *TaskProcessor) Subscribe(fn onStatusUpdateFn) {
	tp.subs = append(tp.subs, fn)
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
