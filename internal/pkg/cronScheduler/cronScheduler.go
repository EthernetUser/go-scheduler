package cronScheduler

import (
	"fmt"
	domainJobs "go-scheduler/internal/domain/jobs"
	"net/http"

	"github.com/robfig/cron/v3"
)

type Scheduler struct{
	c *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		c: cron.New(),
	}
}

func (s *Scheduler) AddJob(job domainJobs.JobForCreation) (int, error) {
	id, err := s.c.AddFunc(job.Cron, func() {
		err := s.doJob(job.Url)
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		return 0, err
	}
	s.c.Start()
	return int(id), nil
}

func (s *Scheduler) DeleteJob(id int) error {
	s.c.Remove(cron.EntryID(id))
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
