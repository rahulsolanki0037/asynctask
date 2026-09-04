package service

import (
	"github.com/rahulsolanki0037/asynctask/internal/model"
	"github.com/rahulsolanki0037/asynctask/internal/queue"
)

type JobRepository interface {
	CreateJob(req model.Job) model.Job
	GetAll() []model.Job
	GetJobById(id int) (model.Job, bool)
	UpdateStatus(id int, status string) bool
}

type JobService struct {
	repository JobRepository
	queue      queue.JobQueue
}

func NewJobService(repository JobRepository, queue queue.JobQueue) *JobService {
	return &JobService{
		repository: repository,
		queue:      queue,
	}
}

func (s *JobService) CreateJob(req model.CreateJob) model.Job {
	job := model.Job{
		Type:    req.Type,
		Payload: req.Payload,
		Status:  "QUEUED",
	}

	job = s.repository.CreateJob(job)

	s.queue.Enqueue(job)

	return job
}

func (s *JobService) GetAllJobs() []model.Job {
	return s.repository.GetAll()
}

func (s *JobService) GetJobById(id int) (model.Job, bool) {
	return s.repository.GetJobById(id)
}

func (s *JobService) UpdateStatus(id int, status string) bool {
	return s.repository.UpdateStatus(id, status)
}
