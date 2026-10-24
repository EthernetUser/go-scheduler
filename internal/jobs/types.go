package jobs

import "time"

type JobForCreation struct {
	Name        string
	Description string
	Cron        string
	Url         string
}

type Job struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cron        string    `json:"cron"`
	Url         string    `json:"url"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
