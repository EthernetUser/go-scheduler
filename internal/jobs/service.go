package jobs

import "context"

type Jobs struct {
	respository Repository
}

func New(repository Repository) *Jobs {
	return &Jobs{
		respository: repository,
	}
}

func (j *Jobs) CreateJob(ctx context.Context, job JobForCreation) error {
	err := j.respository.InsertJob(ctx, job)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jobs) DeleteJob(ctx context.Context, id int) error {
	err := j.respository.DeleteJob(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jobs) GetJobs(ctx context.Context) ([]Job, error) {
	return j.respository.GetJobs(ctx)
}
