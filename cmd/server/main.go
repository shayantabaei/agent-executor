package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/executions", executionHandler)

	fmt.Println("Agent Executor listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := json.Marshal(map[string]string{"status": "ok"})
	if err != nil {
		log.Printf("Error marshaling health response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data = append(data, '\n') // Add a newline for better readability in logs
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing health response: %v", err)
	}

}

type ExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ExecutionResponse struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}

func executionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExecutionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		http.Error(w, "Language is required", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	resp := ExecutionResponse{
		Language: req.Language,
		Code:     req.Code,
		Status:   "accepted",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshaling execution response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data = append(data, '\n') // Add a newline for better readability in logs
	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing execution response: %v", err)
	}
}
