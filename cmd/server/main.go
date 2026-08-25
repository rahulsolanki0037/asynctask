package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Register the HTTP handler for the /health endpoint
	http.HandleFunc("/health", healthHandler)

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