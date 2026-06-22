package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shayantabaei/agent-executor/internal/execution"
)

type Handler struct {
	executor execution.Executor
	config   Config
}

func NewHandler(executor execution.Executor) *Handler {
	return NewHandlerWithConfig(executor, DefaultConfig())
}

func NewHandlerWithConfig(executor execution.Executor, config Config) *Handler {
	return &Handler{executor: executor, config: config}
}

// HealthHandler reports whether the HTTP server is running.
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// ExecutionHandler validates and executes a submitted code request.
func (h *Handler) ExecutionHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit the request body before decoding JSON to avoid reading
	// oversized payloads into memory.
	r.Body = http.MaxBytesReader(w, r.Body, h.config.MaxBodySize)

	var request ExecutionRequest
	// Decode JSON body into ExecutionRequest struct
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if request.Language == "" {
		http.Error(w, "Language is required", http.StatusBadRequest)
		return
	}

	if request.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	if len(request.Code) > h.config.MaxCodeSize {
		http.Error(w, "Code exceeds maximum size", http.StatusBadRequest)
	}

	// Create a 5 second timeout context for the execution request to prevent long-running code.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call the executor to run the code and capture the result
	// Pass the request context so execution can respond to cancellation and timeouts.
	result, err := h.executor.Run(
		ctx,
		execution.Request{
			Language: request.Language,
			Code:     request.Code,
		},
	)

	var unsupportedLanguageError execution.UnsupportedLanguageError
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Execution timed out", http.StatusRequestTimeout)
		}

		if errors.As(err, &unsupportedLanguageError) {
			http.Error(w, unsupportedLanguageError.Error(), http.StatusBadRequest)
		}

		log.Printf("Error executing code: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the execution result back as JSON
	writeJSON(w, http.StatusOK, ExecutionResponse{
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
		ExitCode: result.ExitCode,
	})
}

func (h *Handler) RuntimesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, RuntimesResponse{
		Runtimes: execution.SupportedLanguages(),
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
