package worker

import (
	"context"
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

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case job, ok := <-w.queue.Jobs():
			if !ok {
				return
			}

			_, err := w.service.UpdateStatus(ctx, job.ID, "PROCESSING")
			if err != nil {
				fmt.Println("Failed to update job status: ", err)
				continue
			}

			err = w.ProcessJob(job)
			if err != nil {
				_, err = w.service.UpdateStatus(ctx, job.ID, "FAILED")
				if err != nil {
					fmt.Println("Failed to update job status: ", err)
				}

				if job.RetryCount <= retries {
					retryJob, err := w.service.RetryJob(ctx, job.ID)
					if err != nil {
						fmt.Println("Retry job failed: ", err)
						continue
					}

					w.queue.Enqueue(retryJob)
				}

				continue
			}

			_, err = w.service.UpdateStatus(ctx, job.ID, "COMPLETED")
			if err != nil {
				fmt.Println("Failed to update job status: ", err)
			}

		case <-ctx.Done():
			return
		}

	}
}
