package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shayantabaei/agent-executor/internal/api"
)

func main() {
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/executions", api.ExecutionHandler)

	fmt.Println("Agent Executor listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
