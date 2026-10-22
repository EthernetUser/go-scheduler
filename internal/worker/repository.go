package worker

import (
	"context"
	"go-scheduler/internal/jobs"
)

type Repository interface {
	GetJobs(context.Context) ([]jobs.Job, error)
}
