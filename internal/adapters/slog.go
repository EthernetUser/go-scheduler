package adapters

import (
	handlers "go-scheduler/internal/handlers/jobs"
	"log/slog"
)

type slogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(logger *slog.Logger) *slogAdapter {
	return &slogAdapter{
		logger: logger,
	}
}

func (s *slogAdapter) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}

func (s *slogAdapter) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

func (s *slogAdapter) WithGroup(name string) handlers.Logger {
	return &slogAdapter{
		logger: s.logger.WithGroup(name),
	}
}
