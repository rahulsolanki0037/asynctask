package service

import (
	"github.com/rahulsolanki0037/asynctask/internal/model"
	"github.com/rahulsolanki0037/asynctask/internal/repository"
)

type JobService struct {
	repository *repository.JobRepository
}

func NewJobService (repository *repository.JobRepository) *JobService {
	return &JobService{
		repository: repository,
	}
}

func (s *JobService) CreateJob(req model.CreateJob) model.Job {
	job := model.Job{
		Type: req.Type,
		Payload: req.Payload,
		Status: "QUEUED",
	}

	return s.repository.CreateJob(job)
}

func (s *JobService) GetAllJobs() []model.Job {
	return s.repository.GetAll()
}

func (s *JobService) GetJobById(id int) (model.Job, bool) {
	return s.repository.GetByJobId(id)
}