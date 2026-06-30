package main

import (
	"context"
	"log"
	"os"

	"github.com/shayantabaei/agent-executor/internal/execution"
	"github.com/shayantabaei/agent-executor/internal/mcpserver"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	service := execution.NewServiceWithConfig(
		execution.NewDockerExecutor(),
		execution.DefaultServiceConfig(),
	)

	server := mcpserver.New(service)

	if err := server.RunStdio(context.Background()); err != nil {
		logger.Printf("mcp server failed: %v", err)
		os.Exit(1)
	}
}
