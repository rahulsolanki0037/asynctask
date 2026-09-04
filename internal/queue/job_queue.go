package queue

import "github.com/rahulsolanki0037/asynctask/internal/model"

type JobQueue struct {
	jobs chan model.Job
}

func NewJobQueue(size int) *JobQueue {
	return &JobQueue{
		jobs: make(chan model.Job, size),
	}
}

func (q *JobQueue) Enqueue(job model.Job) {
	q.jobs <- job
}

func (q *JobQueue) Jobs() <- chan model.Job {
	return q.jobs
} 