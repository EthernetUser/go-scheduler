package worker

import (
	"context"
	"go-scheduler/internal/jobs"
)

type Executor interface {
	Execute(ctx context.Context, job jobs.Job) error
}
