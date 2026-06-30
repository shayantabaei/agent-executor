package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shayantabaei/agent-executor/internal/execution"
)

const (
	ServerName    = "agent-executor"
	ServerVersion = "0.1.0"
)

type Server struct {
	service *execution.Service
}

func New(service *execution.Service) *Server {
	return &Server{service: service}
}

func (s *Server) MCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Version: ServerVersion,
	}, nil)

	s.addResources(server)
	s.addTools(server)

	return server
}

func (s *Server) RunStdio(ctx context.Context) error {
	return s.MCPServer().Run(ctx, &mcp.StdioTransport{})
}
