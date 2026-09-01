package worker

import (
	"fmt"

	"github.com/rahulsolanki0037/asynctask/internal/queue"
)

type Worker struct {
	queue *queue.JobQueue
}

func NewWorker (queue *queue.JobQueue) *Worker {
	return &Worker{
		queue: queue,
	}
}

func (w *Worker) Start() {
	for job := range w.queue.Jobs() {
		fmt.Printf("Processing Job for %d\n", job.ID)
	}
}