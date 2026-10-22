package worker

import (
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
