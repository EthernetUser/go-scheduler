package worker

import (
	"context"
	"fmt"
	"go-scheduler/internal/jobs"
	"sync"
	"time"
)

type Worker struct {
	repository Repository
	executor   Executor
	cron       Cron
	interval   time.Duration
	mu         sync.Mutex
	scheduled  map[int]int
	seen       map[int]jobs.Job
}

func New(
	repository Repository,
	executor Executor,
	interval time.Duration,
	cron Cron,
) *Worker {
	return &Worker{
		repository: repository,
		executor:   executor,
		cron:       cron,
		interval:   interval,
		scheduled:  make(map[int]int),
		seen:       make(map[int]jobs.Job),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.cron.Start()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:

		}
	}
}

func (w *Worker) Shutdown(shutdownCtx context.Context) error {
	cronCtx := w.cron.Stop()
	select {
	case <-cronCtx.Done():
		return nil
	case <-shutdownCtx.Done():
		return fmt.Errorf("Timeout occurred")
	}
}
