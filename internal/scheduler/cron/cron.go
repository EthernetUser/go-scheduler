package cron

import (
	"fmt"
	domainJobs "go-scheduler/internal/jobs"
	"net/http"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) AddJob(job domainJobs.JobForCreation) (int, error) {
	id, err := s.cron.AddFunc(job.Cron, func() {
		err := s.doJob(job.Url)
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		return 0, err
	}
	s.cron.Start()
	return int(id), nil
}

func (s *Scheduler) DeleteJob(id int) error {
	s.cron.Remove(cron.EntryID(id))
	return nil
}

func (s *Scheduler) doJob(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(res)
	return nil
}
