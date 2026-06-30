package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayantabaei/agent-executor/internal/api"
	"github.com/shayantabaei/agent-executor/internal/execution"
)

func main() {
	executor := execution.NewDockerExecutor()
	// Inject the executor into the HTTP handler
	handler := api.NewHandler(executor)

	// Wire up the HTTP server
	http.HandleFunc("/health", handler.HealthHandler)
	http.HandleFunc("/executions", handler.ExecutionHandler)
	http.HandleFunc("/runtimes", handler.RuntimesHandler)

	fmt.Println("Agent Executor listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
