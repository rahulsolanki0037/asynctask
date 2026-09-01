package repository

import (
	"sync"

	"github.com/rahulsolanki0037/asynctask/internal/model"
)

type JobRepository struct {
	mu     sync.Mutex
	jobs   map[int]model.Job
	nextID int
}

func NewJobRepository() *JobRepository {
	return &JobRepository{
		jobs:   make(map[int]model.Job),
		nextID: 1,
	}
}

func (r *JobRepository) CreateJob(job model.Job) model.Job {
	r.mu.Lock()
	defer r.mu.Unlock()

	job.ID = r.nextID
	r.nextID++

	r.jobs[job.ID] = job

	return job
}

func (r *JobRepository) GetAll() []model.Job {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobs := make([]model.Job, 0, len(r.jobs))

	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

func (r *JobRepository) GetJobById(jobId int) (model.Job, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	job, exists := r.jobs[jobId]
	return job, exists
}

func (r *JobRepository) UpdateStatus(jobId int, status string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[jobId]
	if !exists {
		return false
	}

	job.Status = status

	r.jobs[jobId] = job

	return true
}