package jobs

type Db interface {
	InsertJob(JobForCreation) error
	GetJobs() ([]Job, error)
	GetJob(int) (Job, error)
	UpdateJob(Job) error
	DeleteJob(int) error
}

type Scheduler interface {
	AddJob(JobForCreation) (int, error)
	DeleteJob(int) error
}

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

type Jobs struct {
	db        Db
	scheduler Scheduler
}

func New(db Db, scheduler Scheduler) *Jobs {
	return &Jobs{
		db:        db,
		scheduler: scheduler,
	}
}

func (j *Jobs) CreateJob(job JobForCreation) (int, error) {
	id, err := j.scheduler.AddJob(job)
	if err != nil {
		return 0, err
	}

	err = j.db.InsertJob(job)
	if err != nil {
		j.scheduler.DeleteJob(id)
		return 0, err
	}

	return id, nil
}

func (j *Jobs) DeleteJob(id int) error {
	job, err := j.db.GetJob(id)
	if err != nil {
		return err
	}

	err = j.scheduler.DeleteJob(job.CronId)
	if err != nil {
		return err
	}
	return j.db.DeleteJob(id)
}

func (j *Jobs) GetJobs() ([]Job, error) {
	return j.db.GetJobs()
}
