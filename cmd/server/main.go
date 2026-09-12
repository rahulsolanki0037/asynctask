package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rahulsolanki0037/asynctask/config"
	"github.com/rahulsolanki0037/asynctask/database"
	"github.com/rahulsolanki0037/asynctask/internal/handler"
	"github.com/rahulsolanki0037/asynctask/internal/queue"
	"github.com/rahulsolanki0037/asynctask/internal/repository"
	"github.com/rahulsolanki0037/asynctask/internal/service"
	"github.com/rahulsolanki0037/asynctask/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}

	dbPool, err := database.Connect(ctx, config.LoadDBConfig().DBUrl())
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer dbPool.Close()

	jobRepository := repository.NewPostgresJobRepository(dbPool)
	jobQueue := queue.NewJobQueue(10)
	jobService := service.NewJobService(jobRepository, *jobQueue)
	jobHandler := handler.NewJobHandler(jobService)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM) // SIGTERM - Signal Terminate

	totalWorkers := 5

	// Worker pool with multiple workers
	for i := 0; i < totalWorkers; i++ {
		jobWorker := worker.NewWorker(i, jobQueue, jobService)
		go jobWorker.Start(ctx)
	}

	recoveredJobs, err := jobService.RecoverProcessingJobs(ctx)
	if err != nil {
		log.Fatal("Failed to recover processing jobs: ", err)
	}

	for _, job := range recoveredJobs {
		jobQueue.Enqueue(job)
	}

	// Register the HTTP handlers
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/createJob", jobHandler.CreateJobHandler)
	http.HandleFunc("/getJobs", jobHandler.GetAllJobs)
	http.HandleFunc("/getJob/", jobHandler.GetByJobId)

	// Create the HTTP Server
	server := http.Server{Addr: ":8080"}

	// Start the HTTP Server
	go func() {
		fmt.Println("Server running on : 8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error: ", err)
		}
	}()

	<-signalChan

	fmt.Println("Server shutdown signal received")

	cancel()

	err = server.Shutdown(context.Background())
	if err != nil {
		fmt.Println("Server shutdown error: ", err)
	}
}

// healthHander handles request made to /health endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Health endpoint only supports GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Set the successful HTTP status code
	w.WriteHeader(http.StatusOK)

	// Write the JSON response body
	fmt.Fprintln(w, `{"status": "UP"}`)
}
