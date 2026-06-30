package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shayantabaei/agent-executor/internal/execution"
)

type executeCodeInput struct {
	Runtime string                `json:"runtime" jsonschema:"runtime to execute, for example python or javascript"`
	Code    string                `json:"code" jsonschema:"source code to execute"`
	Files   []execution.InputFile `json:"files,omitempty" jsonschema:"optional files to make available during execution"`
}

type executeCodeOutput struct {
	Stdout    string               `json:"stdout"`
	Stderr    string               `json:"stderr"`
	ExitCode  int                  `json:"exitCode"`
	Artifacts []execution.Artifact `json:"artifacts,omitempty"`
}

func (s *Server) addTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "execute_code",
		Description: "Execute code in a supported docker backed runtime and return stdout, stderr, exit code, timeout status and artifacts",
	}, s.executeCode)
}

func (s *Server) executeCode(ctx context.Context, request *mcp.CallToolRequest, input executeCodeInput) (*mcp.CallToolResult, executeCodeOutput, error) {
	result, err := s.service.Run(ctx, execution.Request{
		Language: input.Runtime,
		Code:     input.Code,
		Files:    input.Files,
	})

	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
		}, executeCodeOutput{}, nil
	}

	output := executeCodeOutput{
		Stdout:    result.Stdout,
		Stderr:    result.Stderr,
		ExitCode:  result.ExitCode,
		Artifacts: result.Artifacts,
	}

	text, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return nil, executeCodeOutput{}, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(text)},
		},
	}, output, nil
}
