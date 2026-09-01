package main

import (
	"fmt"
	"net/http"

	"github.com/rahulsolanki0037/asynctask/internal/handler"
	"github.com/rahulsolanki0037/asynctask/internal/queue"
	"github.com/rahulsolanki0037/asynctask/internal/repository"
	"github.com/rahulsolanki0037/asynctask/internal/service"
)

func main() {
	jobRepository := repository.NewJobRepository()
	jobQueue := queue.NewJobQueue(10)
	jobService := service.NewJobService(jobRepository, *jobQueue)
	jobHandler := handler.NewJobHandler(jobService)

	// Register the HTTP handlers
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/createJob", jobHandler.CreateJobHandler)
	http.HandleFunc("/getJobs", jobHandler.GetAllJobs)
	http.HandleFunc("/getJob/", jobHandler.GetByJobId)

	// Start the HTTP Server on port 8080
	fmt.Println("Server running on : 8080")
	http.ListenAndServe(":8080", nil)
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
