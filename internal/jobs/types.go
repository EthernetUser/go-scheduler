package jobs

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
	CronId      int
	Url         string
}
