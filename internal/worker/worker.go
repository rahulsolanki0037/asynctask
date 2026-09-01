package worker

import (
	"fmt"
	"time"

	"github.com/rahulsolanki0037/asynctask/internal/queue"
	"github.com/rahulsolanki0037/asynctask/internal/service"
)

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

func (w *Worker) Start() {
	for job := range w.queue.Jobs() {
		fmt.Printf("Worker %d processing Job for %d\n", w.id, job.ID)

		w.service.UpdateStatus(job.ID, "PROCESSING")

		time.Sleep(3 * time.Second)

		w.service.UpdateStatus(job.ID, "COMPLETED")
	}
}
