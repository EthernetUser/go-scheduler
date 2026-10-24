package handlers

import "go-scheduler/internal/jobs"

type Logger interface {
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	WithGroup(name string) Logger
}

type handlers struct {
	logger     Logger
	jobService *jobs.Jobs
}

func New(logger Logger, jobService *jobs.Jobs) *handlers {
	return &handlers{
		logger:     logger,
		jobService: jobService,
	}
}
