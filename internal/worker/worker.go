package worker

import (
	"fmt"

	"github.com/rahulsolanki0037/asynctask/internal/queue"
)

type Worker struct {
	id int
	queue *queue.JobQueue
}

func NewWorker (id int, queue *queue.JobQueue) *Worker {
	return &Worker{
		id: id,
		queue: queue,
	}
}

func (w *Worker) Start() {
	for job := range w.queue.Jobs() {
		fmt.Printf("Worker %d processing Job for %d\n", w.id, job.ID)
	}
}