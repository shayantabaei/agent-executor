// internal/mcpserver/server_test.go
package mcpserver

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestListsResources(t *testing.T) {
	ctx := context.Background()

	server := New(nil).MCPServer()
	session := connectTestClient(t, ctx, server)
	defer session.Close()

	result, err := session.ListResources(ctx, &mcp.ListResourcesParams{})
	if err != nil {
		t.Fatalf("ListResources returned error: %v", err)
	}

	got := map[string]bool{}
	for _, resource := range result.Resources {
		got[resource.URI] = true
	}

	if !got[RuntimesURI] {
		t.Fatalf("expected runtimes resource to be listed")
	}

	if !got[CapabilitiesURI] {
		t.Fatalf("expected capabilities resource to be listed")
	}
}

func TestReadsRuntimesResource(t *testing.T) {
	ctx := context.Background()

	server := New(nil).MCPServer()
	session := connectTestClient(t, ctx, server)
	defer session.Close()

	result, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: RuntimesURI,
	})
	if err != nil {
		t.Fatalf("ReadResource returned error: %v", err)
	}

	if len(result.Contents) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Contents))
	}

	content := result.Contents[0]
	if content.URI != RuntimesURI {
		t.Fatalf("expected URI %q, got %q", RuntimesURI, content.URI)
	}

	if content.MIMEType != "application/json" {
		t.Fatalf("expected application/json, got %q", content.MIMEType)
	}

	if content.Text == "" {
		t.Fatalf("expected resource text")
	}
}

func TestReadsCapabilitiesResource(t *testing.T) {
	ctx := context.Background()

	server := New(nil).MCPServer()
	session := connectTestClient(t, ctx, server)
	defer session.Close()

	result, err := session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: CapabilitiesURI,
	})
	if err != nil {
		t.Fatalf("ReadResource returned error: %v", err)
	}

	if len(result.Contents) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Contents))
	}

	content := result.Contents[0]
	if content.URI != CapabilitiesURI {
		t.Fatalf("expected URI %q, got %q", CapabilitiesURI, content.URI)
	}

	if content.MIMEType != "application/json" {
		t.Fatalf("expected application/json, got %q", content.MIMEType)
	}

	if content.Text == "" {
		t.Fatalf("expected resource text")
	}
}

func connectTestClient(t *testing.T, ctx context.Context, server *mcp.Server) *mcp.ClientSession {
	t.Helper()

	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("server connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "agent-executor-test-client",
		Version: "0.1.0",
	}, nil)

	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}

	return session
}
