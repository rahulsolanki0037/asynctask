package service

import (
	"context"

	"github.com/rahulsolanki0037/asynctask/internal/model"
	"github.com/rahulsolanki0037/asynctask/internal/queue"
)

type JobRepository interface {
	CreateJob(ctx context.Context, req model.Job) (model.Job, error)
	GetAll(ctx context.Context) ([]model.Job, error)
	GetJobById(ctx context.Context, id int) (model.Job, error)
	UpdateStatus(ctx context.Context, id int, status string) (bool, error)
	RetryJob(ctx context.Context, id int) (model.Job, error)
	RecoverProcessingJobs(ctx context.Context) ([]model.Job, error)
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

func (s *JobService) CreateJob(ctx context.Context, req model.CreateJob) (model.Job, error) {
	job := model.Job{
		Type:    req.Type,
		Payload: req.Payload,
		Status:  "QUEUED",
	}

	job, err := s.repository.CreateJob(ctx, job)
	if err != nil {
		return model.Job{}, err
	}

	s.queue.Enqueue(job)

	return job, nil
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]model.Job, error) {
	jobs, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *JobService) GetJobById(ctx context.Context, id int) (model.Job, error) {
	job, err := s.repository.GetJobById(ctx, id)
	if err != nil {
		return model.Job{}, err
	}

	return job, nil
}

func (s *JobService) UpdateStatus(ctx context.Context, id int, status string) (bool, error) {
	result, err := s.repository.UpdateStatus(ctx, id, status)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (s *JobService) RetryJob(ctx context.Context, id int) (model.Job, error) {
	job, err := s.repository.RetryJob(ctx, id)
	if err != nil {
		return model.Job{}, err
	}

	return job, nil

}

func (s *JobService) RecoverProcessingJobs(ctx context.Context) ([]model.Job, error) {
	jobs, err := s.repository.RecoverProcessingJobs(ctx)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}