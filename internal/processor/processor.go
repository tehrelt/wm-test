package processor

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tehrelt/wm-test/internal/models"
)

type StatusUpdateFn func(context.Context, *models.Task) error

type processingTask struct {
	*models.Task
	update StatusUpdateFn
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

type TaskProcessor struct {
	queue   chan *processingTask
	factory *processingTaskFactory
	logger  *slog.Logger
}

func NewTaskProcessor(update StatusUpdateFn, queueSize int) *TaskProcessor {
	tp := &TaskProcessor{
		queue:   make(chan *processingTask, queueSize),
		factory: &processingTaskFactory{update: update},
		logger:  slog.With(slog.String("comp", "processor.TaskProcessor")),
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
	tp.queue <- tp.factory.New(task)
}
