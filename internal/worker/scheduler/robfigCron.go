package scheduler

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type RobfigCron struct {
	cron *cron.Cron
}

func NewRobfigCronScheduler() *RobfigCron {
	return &RobfigCron{
		cron: cron.New(),
	}
}

func (s *RobfigCron) AddFunc(spec string, cmd func()) (int, error) {
	id, err := s.cron.AddFunc(spec, cmd)
	if err != nil {
		return 0, fmt.Errorf("error occurred on adding cron func: %w", err)
	}
	return int(id), nil
}

func (s *RobfigCron) Remove(id int) {
	s.cron.Remove(cron.EntryID(id))
}

func (s *RobfigCron) Start() {
	s.cron.Start()
}

func (s *RobfigCron) Stop() context.Context {
	return s.cron.Stop()
}
