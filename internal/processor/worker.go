package processor

import (
	"context"
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
		logger: slog.With(slog.String("comp", "processor.worker"), slog.Int("id", id)),
	}
}

func (w *worker) run(ctx context.Context) {
	for task := range w.queue {
		task.SetStatus(ctx, models.PSProcessing)

		min := 3
		max := 10
		sec := rand.IntN(max-min) + min

		w.logger.Info("processing task", "task", task.Id, "sec", sec)
		time.Sleep(time.Duration(sec) * time.Second)

		task.SetStatus(ctx, models.PSDone)
		w.logger.Info("task done", "task", task.Id)
	}
}
