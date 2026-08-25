package repository

import "github.com/rahulsolanki0037/asynctask/internal/model"

type JobRepository struct {
	jobs map[int]model.Job
	nextID int
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		jobs: make(map[int]model.Job),
		nextID: 1,
	}
}

func (r *JobRepository) CreateJob(job model.Job) model.Job {
	job.ID = r.nextID
	r.nextID++

	r.jobs[job.ID] = job

	return job
}

func (r *JobRepository) GetAll() []model.Job {
	jobs := make([]model.Job, 0, len(r.jobs))

	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

func (r *JobRepository) GetByJobId(jobId int) (model.Job,bool) {
	job, exists := r.jobs[jobId]
	return job, exists
}