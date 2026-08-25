package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rahulsolanki0037/asynctask/internal/model"
	"github.com/rahulsolanki0037/asynctask/internal/repository"
)

type JobHandler struct {
	repository *repository.JobRepository
}

func NewJobHandler(repository *repository.JobRepository) *JobHandler {
	return &JobHandler{
		repository: repository,
	}
}

func (h *JobHandler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	// Only POST requests are allowed for creating a job.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateJob
	// Decode the JSON request body into the Go struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if (err != nil) {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job := model.Job{
		Type: req.Type,
		Payload: req.Payload,
		Status: "Queued",
	}

	job = h.repository.CreateJob(job)

	// Tell the client that a new resource was successfully created.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Convert the Job struct into JSON and send it to the client.
	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobs := h.repository.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetByJobId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "getJob" {
		http.Error(w, "Invalid job path", http.StatusBadRequest)
		return
	}

	jobId, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid Job Id", http.StatusBadRequest)
		return
	}

	job, exists := h.repository.GetByJobId(jobId)
	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(job)
}