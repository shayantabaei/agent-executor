package execution

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestWorkspaceCollectArtifactsSkipsDirectories(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	if err := os.MkdirAll(filepath.Join(ws.path, "nested"), 0755); err != nil {
		t.Fatalf("unexpected mkdir error: %v", err)
	}

	artifacts, err := ws.collectArtifacts(nil, DefaultDockerConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(artifacts) != 0 {
		t.Fatalf("expected no artifacts, got %d: %+v", len(artifacts), artifacts)
	}
}

func TestWorkspaceCollectArtifactsRejectsTooManyArtifacts(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	for i := 0; i < 2; i++ {
		path := filepath.Join(ws.path, fmt.Sprintf("output-%d.txt", i))
		if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
			t.Fatalf("unexpected write error: %v", err)
		}
	}

	cfg := DefaultDockerConfig()
	cfg.MaxArtifactCount = 1

	_, err = ws.collectArtifacts(nil, cfg)
	if !errors.Is(err, ErrTooManyArtifacts) {
		t.Fatalf("expected ErrTooManyArtifacts, got %v", err)
	}
}
func TestWorkspaceCollectArtifactsRejectsLargeArtifact(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	if err := os.WriteFile(filepath.Join(ws.path, "large.txt"), []byte("abcd"), 0644); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	cfg := DefaultDockerConfig()
	cfg.MaxArtifactSizeBytes = 3

	_, err = ws.collectArtifacts(nil, cfg)
	if !errors.Is(err, ErrArtifactTooLarge) {
		t.Fatalf("expected ErrArtifactTooLarge, got %v", err)
	}
}
func TestWorkspaceCollectArtifactsRejectsLargeTotalArtifactSize(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	if err := os.WriteFile(filepath.Join(ws.path, "a.txt"), []byte("abc"), 0644); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	if err := os.WriteFile(filepath.Join(ws.path, "b.txt"), []byte("def"), 0644); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	cfg := DefaultDockerConfig()
	cfg.MaxArtifactSizeBytes = 10
	cfg.MaxTotalArtifactSizeBytes = 5

	_, err = ws.collectArtifacts(nil, cfg)
	if !errors.Is(err, ErrTotalArtifactsTooLarge) {
		t.Fatalf("expected ErrTotalArtifactsTooLarge, got %v", err)
	}
}
func TestWorkspaceCollectArtifactsDoesNotInlineLargeArtifact(t *testing.T) {
	ws, err := createWorkspace(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ws.cleanup()

	if err := os.WriteFile(filepath.Join(ws.path, "large.txt"), []byte("abcd"), 0644); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	cfg := DefaultDockerConfig()
	cfg.MaxInlineArtifactBytes = 3

	artifacts, err := ws.collectArtifacts(nil, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(artifacts) != 1 {
		t.Fatalf("expected 1 artifact, got %d", len(artifacts))
	}

	if artifacts[0].Content != "" {
		t.Fatalf("expected large artifact content not to be inlined, got %q", artifacts[0].Content)
	}

	if artifacts[0].Encoding != "" {
		t.Fatalf("expected large artifact encoding to be empty, got %q", artifacts[0].Encoding)
	}
}
func TestDockerExecutorReturnsGeneratedArtifact(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `open("output.txt", "w").write("hello artifact")`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Artifacts) != 1 {
		t.Fatalf("expected 1 artifact, got %d: %+v", len(result.Artifacts), result.Artifacts)
	}

	artifact := result.Artifacts[0]

	if artifact.Path != "output.txt" {
		t.Fatalf("expected artifact path %q, got %q", "output.txt", artifact.Path)
	}

	if artifact.Size != int64(len("hello artifact")) {
		t.Fatalf("expected artifact size %d, got %d", len("hello artifact"), artifact.Size)
	}

	if artifact.Content != "hello artifact" {
		t.Fatalf("expected artifact content %q, got %q", "hello artifact", artifact.Content)
	}

	if artifact.Encoding != "utf-8" {
		t.Fatalf("expected artifact encoding %q, got %q", "utf-8", artifact.Encoding)
	}

	if artifact.ContentType == "" {
		t.Fatal("expected artifact content type to be set")
	}
}
func TestDockerExecutorExcludesInputFilesFromArtifacts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `open("output.txt", "w").write(open("input.txt").read().upper())`,
		Files: []InputFile{
			{
				Path:    "input.txt",
				Content: "hello",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Artifacts) != 1 {
		t.Fatalf("expected 1 artifact, got %d: %+v", len(result.Artifacts), result.Artifacts)
	}

	artifact := result.Artifacts[0]

	if artifact.Path != "output.txt" {
		t.Fatalf("expected artifact path %q, got %q", "output.txt", artifact.Path)
	}

	if artifact.Content != "HELLO" {
		t.Fatalf("expected artifact content %q, got %q", "HELLO", artifact.Content)
	}
}

func TestDockerExecutorCleansUpWorkspaceAfterRuntimeError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	realCreateWorkspace := createWorkspace

	var cleaned bool

	executor := NewDockerExecutor()
	executor.createWorkspace = realCreateWorkspace
	executor.cleanupWorkspace = func(ws *workspace) {
		cleaned = true
		if err := ws.cleanup(); err != nil {
			t.Fatalf("unexpected cleanup error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `raise Exception("boom")`,
	})

	if err != nil {
		t.Fatalf("expected runtime error to be returned as result, got err: %v", err)
	}

	if result.ExitCode == 0 {
		t.Fatal("expected non-zero exit code")
	}

	if !cleaned {
		t.Fatal("expected workspace cleanup to be called")
	}
}

func TestDockerExecutorCleansUpWorkspaceAfterArtifactError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	var cleaned bool

	cfg := DefaultDockerConfig()
	cfg.MaxArtifactCount = 0

	executor := NewDockerExecutorWithConfig(cfg)
	executor.cleanupWorkspace = func(ws *workspace) {
		cleaned = true
		if err := ws.cleanup(); err != nil {
			t.Fatalf("unexpected cleanup error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `open("output.txt", "w").write("hello")`,
	})

	if !errors.Is(err, ErrTooManyArtifacts) {
		t.Fatalf("expected ErrTooManyArtifacts, got %v", err)
	}

	if !cleaned {
		t.Fatal("expected workspace cleanup to be called")
	}
}

func TestDockerExecutorCleansUpWorkspaceAfterSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	var cleaned bool

	executor := NewDockerExecutor()
	executor.cleanupWorkspace = func(ws *workspace) {
		cleaned = true
		if err := ws.cleanup(); err != nil {
			t.Fatalf("unexpected cleanup error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `print("ok")`,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cleaned {
		t.Fatal("expected workspace cleanup to be called")
	}
}

func TestDockerExecutorCleansUpWorkspaceAfterTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	var cleaned bool

	executor := NewDockerExecutor()
	executor.cleanupWorkspace = func(ws *workspace) {
		cleaned = true
		if err := ws.cleanup(); err != nil {
			t.Fatalf("unexpected cleanup error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `import time; time.sleep(5)`,
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}

	if !cleaned {
		t.Fatal("expected workspace cleanup to be called")
	}
}
