package robfigcron

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type robfigCron struct {
	cron *cron.Cron
}

func New() *robfigCron {
	return &robfigCron{
		cron: cron.New(),
	}
}

func (s *robfigCron) AddFunc(spec string, cmd func()) (int, error) {
	id, err := s.cron.AddFunc(spec, cmd)
	if err != nil {
		return 0, fmt.Errorf("error occurred on adding cron func: %w", err)
	}
	return int(id), nil
}

func (s *robfigCron) Remove(id int) {
	s.cron.Remove(cron.EntryID(id))
}

func (s *robfigCron) Start() {
	s.cron.Start()
}

func (s *robfigCron) Stop() context.Context {
	return s.cron.Stop()
}
