package jobs

import "context"

type Scheduler interface {
	AddJob(JobForCreation) (int, error)
	DeleteJob(int) error
}

type Jobs struct {
	respository Repository
	scheduler   Scheduler
}

func New(repository Repository, scheduler Scheduler) *Jobs {
	return &Jobs{
		respository: repository,
		scheduler:   scheduler,
	}
}

func (j *Jobs) CreateJob(ctx context.Context, job JobForCreation) (int, error) {
	id, err := j.scheduler.AddJob(job)
	if err != nil {
		return 0, err
	}

	err = j.respository.InsertJob(ctx, job)
	if err != nil {
		j.scheduler.DeleteJob(id)
		return 0, err
	}

	return id, nil
}

func (j *Jobs) DeleteJob(ctx context.Context, id int) error {
	job, err := j.respository.GetJob(ctx, id)
	if err != nil {
		return err
	}

	err = j.scheduler.DeleteJob(job.CronId)
	if err != nil {
		return err
	}
	return j.respository.DeleteJob(ctx, id)
}

func (j *Jobs) GetJobs(ctx context.Context) ([]Job, error) {
	return j.respository.GetJobs(ctx)
}
