package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayantabaei/agent-executor/internal/api"
	"github.com/shayantabaei/agent-executor/internal/execution"
)

func main() {
	// Use a fake executor until Docker execution is implemented
	executor := execution.FakeExecutor{}
	// Inject the executor into the HTTP handler
	handler := api.NewHandler(executor)

	// Wire up the HTTP server
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/executions", handler.ExecutionHandler)

	fmt.Println("Agent Executor listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
