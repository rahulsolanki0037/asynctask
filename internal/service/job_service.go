package service

import (
	"github.com/rahulsolanki0037/asynctask/internal/model"
)

type JobRepository interface {
	CreateJob(req model.Job) model.Job
	GetAll() []model.Job
	GetJobById(id int) (model.Job, bool)
}

type JobService struct {
	repository JobRepository
}

func NewJobService(repository JobRepository) *JobService {
	return &JobService{
		repository: repository,
	}
}

func (s *JobService) CreateJob(req model.CreateJob) model.Job {
	job := model.Job{
		Type:    req.Type,
		Payload: req.Payload,
		Status:  "QUEUED",
	}

	return s.repository.CreateJob(job)
}

func (s *JobService) GetAllJobs() []model.Job {
	return s.repository.GetAll()
}

func (s *JobService) GetJobById(id int) (model.Job, bool) {
	return s.repository.GetJobById(id)
}
