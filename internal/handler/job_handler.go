package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rahulsolanki0037/asynctask/internal/model"
	"github.com/rahulsolanki0037/asynctask/internal/service"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{
		service: service,
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
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	job, err := h.service.CreateJob(r.Context(), req)
	if err != nil {
		error := fmt.Sprintf("Failed to create job : %s", err)
		http.Error(w, error, http.StatusInternalServerError)
		return
	}

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

	jobs, err := h.service.GetAllJobs(r.Context())
	if err != nil {
		error := fmt.Sprintf("Failed to fetch job details: %s", err)
		http.Error(w, error, http.StatusInternalServerError)
		return
	}

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

	job, err := h.service.GetJobById(r.Context(), jobId)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(job)
}
