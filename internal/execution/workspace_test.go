package execution

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateWorkspaceCreatesTemporaryDirectory(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	info, err := os.Stat(ws.path)
	if err != nil {
		t.Fatalf("expected workspace directory to exist: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("expected workspace path to be a directory")
	}
}

func TestCreateWorkspaceWritesInputFiles(t *testing.T) {
	ws, err := createWorkspace([]InputFile{
		{
			Path:    "input.txt",
			Content: "hello",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	data, err := os.ReadFile(filepath.Join(ws.path, "input.txt"))
	if err != nil {
		t.Fatalf("expected input file to be written: %v", err)
	}

	if string(data) != "hello" {
		t.Fatalf("expected file content %q, got %q", "hello", string(data))
	}
}

func TestCreateWorkspaceWritesNestedInputFiles(t *testing.T) {
	ws, err := createWorkspace([]InputFile{
		{
			Path:    "data/input.txt",
			Content: "nested hello",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	data, err := os.ReadFile(filepath.Join(ws.path, "data", "input.txt"))
	if err != nil {
		t.Fatalf("expected nested input file to be written: %v", err)
	}

	if string(data) != "nested hello" {
		t.Fatalf("expected file content %q, got %q", "nested hello", string(data))
	}
}

func TestWorkspaceCleanupRemovesDirectory(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := ws.path

	if err := ws.cleanup(); err != nil {
		t.Fatalf("unexpected cleanup error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected workspace directory to be removed, stat err: %v", err)
	}
}
