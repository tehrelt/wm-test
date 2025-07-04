package processor

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/tehrelt/wm-test/internal/models"
)

type worker struct {
	id     int
	queue  chan *processingTask
	logger *slog.Logger
}

func newWorker(id int, queue chan *processingTask) *worker {
	return &worker{
		id:     id,
		queue:  queue,
		logger: slog.With(slog.String("comp", fmt.Sprintf("processor.worker/%d", id))),
	}
}

func (w *worker) run(ctx context.Context) {
	for task := range w.queue {
		task.SetStatus(ctx, models.PSProcessing)

		min := 3
		max := 10
		sec := rand.IntN(max-min) + min

		w.logger.Info("processing task", slog.String("task", task.Id.String()), slog.Int("sec", sec))

		select {
		case <-task.ctx.Done():
			w.logger.Info("task canceled", slog.String("task", task.Id.String()))
			task.SetStatus(ctx, models.PSError)
			continue
		case <-time.After(time.Duration(sec) * time.Second):
		}

		task.SetStatus(ctx, models.PSDone)
		w.logger.Info("task done", slog.String("task", task.Id.String()))
	}
}
