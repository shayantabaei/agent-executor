package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shayantabaei/agent-executor/internal/execution"
)

const (
	RuntimesURI     = "agent-executor://runtimes"
	CapabilitiesURI = "agent-executor://capabilities"
)

type runtimeInfo struct {
	Name string `json:"name"`
}

type capabilitiesInfo struct {
	Tools     []string `json:"tools"`
	Resources []string `json:"resources"`
}

func (s *Server) addResources(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		URI:         RuntimesURI,
		Name:        "runtimes",
		Title:       "Supported runtimes",
		Description: "Lists code runtimes supported by agent-executor",
		MIMEType:    "application/json",
	}, s.readRuntimes)

	server.AddResource(&mcp.Resource{
		URI:         CapabilitiesURI,
		Name:        "capabilities",
		Title:       "Agent Executor capabilities",
		Description: "Describes the MCP tools and resources exposed by agent-executor",
		MIMEType:    "application/json",
	}, s.readCapabilities)
}

func (s *Server) readRuntimes(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	_ = ctx

	if req.Params.URI != RuntimesURI {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}
	supportedLanguages := execution.SupportedLanguages()
	runtimes := make([]runtimeInfo, 0, len(supportedLanguages))
	for _, language := range supportedLanguages {
		runtimes = append(runtimes, runtimeInfo{Name: language})
	}

	return jsonResource(RuntimesURI, runtimes)
}

func (s *Server) readCapabilities(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	_ = ctx

	if req.Params.URI != CapabilitiesURI {
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	capabilities := capabilitiesInfo{
		Tools: []string{
			"execute_code",
		},
		Resources: []string{
			RuntimesURI,
			CapabilitiesURI,
		},
	}

	return jsonResource(CapabilitiesURI, capabilities)
}

func jsonResource(uri string, value any) (*mcp.ReadResourceResult, error) {
	data, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(data),
			},
		},
	}, nil
}
