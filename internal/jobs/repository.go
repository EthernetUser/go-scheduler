package jobs

import "context"

type Repository interface {
	InsertJob(context.Context, JobForCreation) error
	GetJobs(context.Context) ([]Job, error)
	GetJob(context.Context, int) (Job, error)
	UpdateJob(context.Context, Job) error
	DeleteJob(context.Context, int) error
}
