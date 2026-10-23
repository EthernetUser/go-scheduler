package jobs

import "time"

type JobForCreation struct {
	Name        string
	Description string
	Cron        string
	Url         string
}

type Job struct {
	Id          int
	Name        string
	Description string
	Cron        string
	Url         string
	UpdatedAt   time.Time
}
