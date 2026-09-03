package worker

import (
	"errors"
	"fmt"
	"time"

	"github.com/rahulsolanki0037/asynctask/internal/model"
	"github.com/rahulsolanki0037/asynctask/internal/queue"
	"github.com/rahulsolanki0037/asynctask/internal/service"
)

const retries = 3

type Worker struct {
	id      int
	queue   *queue.JobQueue
	service *service.JobService
}

func NewWorker(id int, queue *queue.JobQueue, service *service.JobService) *Worker {
	return &Worker{
		id:      id,
		queue:   queue,
		service: service,
	}
}

func (w *Worker) ProcessJob(job model.Job) error {
	fmt.Printf("Worker %d processing Job for %d\n", w.id, job.ID)

	time.Sleep(3 * time.Second)

	if job.Type == "FAIL" {
		return errors.New("Job Processing failed")
	}

	return nil
}

func (w *Worker) Start() {
	for job := range w.queue.Jobs() {
		w.service.UpdateStatus(job.ID, "PROCESSING")

		err := w.ProcessJob(job)
		if err != nil {
			w.service.UpdateStatus(job.ID, "FAILED")
			
			if job.RetryCount < retries {
				retryJob, ok := w.service.RetryJob(job.ID)
				if ok {
					w.queue.Enqueue(retryJob)
				}
			}
			
			continue
		}

		w.service.UpdateStatus(job.ID, "COMPLETED")
	}
}
