package worker

import (
	"context"
	"fmt"
	"go-scheduler/internal/jobs"
	"sync"
	"time"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Worker struct {
	repository Repository
	logger     Logger
	executor   Executor
	cron       Cron
	interval   time.Duration
	jobTimeout time.Duration
	mu         sync.Mutex
	scheduled  map[int]int
	seen       map[int]jobs.Job
}

func New(
	repository Repository,
	logger Logger,
	executor Executor,
	interval time.Duration,
	cron Cron,
	jobTimeout time.Duration,
) *Worker {
	return &Worker{
		repository: repository,
		logger:     logger,
		executor:   executor,
		cron:       cron,
		interval:   interval,
		jobTimeout: jobTimeout,
		scheduled:  make(map[int]int),
		seen:       make(map[int]jobs.Job),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.cron.Start()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.reconcile(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.reconcile(ctx)
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

func (w *Worker) reconcile(ctx context.Context) {
	dbJobs, err := w.repository.GetJobs(ctx)
	if err != nil {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	wanted := make(map[int]jobs.Job, len(dbJobs))
	for _, job := range dbJobs {
		wanted[job.Id] = job
	}

	for id, entryId := range w.scheduled {
		if _, ok := wanted[id]; !ok {
			w.cron.Remove(entryId)
			delete(w.scheduled, id)
			delete(w.seen, id)
		}
	}

	for id, job := range wanted {
		oldJob, exists := w.seen[id]

		if exists && oldJob.UpdatedAt.Equal(job.UpdatedAt) {
			continue
		}

		if exists {
			w.cron.Remove(w.scheduled[id])
			delete(w.scheduled, id)
		}

		jobCopy := job
		entryId, err := w.cron.AddFunc(job.Cron, func() {
			w.executeJob(jobCopy)
		})

		if err != nil {
			w.logger.Error("Failed to add cron", "jobId", jobCopy.Id)
			continue
		}

		w.scheduled[id] = entryId
		w.seen[id] = jobCopy
	}
}

func (w *Worker) executeJob(job jobs.Job) {
	ctx, cancel := context.WithTimeout(context.Background(), w.jobTimeout)
	defer cancel()

	start := time.Now()

	err := w.executor.Execute(ctx, job)

	duration := time.Since(start)
	if err != nil {
		w.logger.Error(
			"Error occurred on processing job",
			"error",
			err,
			"duration",
			duration,
			"jobId",
			job.Id,
		)
		return
	}

	w.logger.Debug("Job has been finished successfully", "duration", duration)
}
