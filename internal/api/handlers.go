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

	if err := validateExecutionRequest(request, h.config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a 5 second timeout context for the execution request to prevent long-running code.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Call the executor to run the code and capture the result
	// Pass the request context so execution can respond to cancellation and timeouts.

	result, err := h.executor.Run(ctx, toExecutionRequest(request))

	var unsupportedLanguageError execution.UnsupportedLanguageError
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Execution timed out", http.StatusRequestTimeout)
			return
		}

		if errors.As(err, &unsupportedLanguageError) {
			http.Error(w, unsupportedLanguageError.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Error executing code: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the execution result back as JSON
	writeJSON(w, http.StatusOK, toExecutionResponse(result))
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

func toExecutionRequest(req ExecutionRequest) execution.Request {
	files := make([]execution.InputFile, 0, len(req.Files))

	for _, file := range req.Files {
		files = append(files, execution.InputFile{
			Path:    file.Path,
			Content: file.Content,
		})
	}

	return execution.Request{
		Language: req.Language,
		Code:     req.Code,
		Files:    files,
	}
}

func toExecutionResponse(result execution.Result) ExecutionResponse {
	artifacts := make([]ArtifactResponse, 0, len(result.Artifacts))

	for _, artifact := range result.Artifacts {
		artifacts = append(artifacts, ArtifactResponse{
			Path:        artifact.Path,
			Size:        artifact.Size,
			Content:     artifact.Content,
			Encoding:    artifact.Encoding,
			ContentType: artifact.ContentType,
		})
	}

	return ExecutionResponse{
		Stdout:    result.Stdout,
		Stderr:    result.Stderr,
		ExitCode:  result.ExitCode,
		Artifacts: artifacts,
	}
}
