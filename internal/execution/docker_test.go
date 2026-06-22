package execution

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDockerExecutorRunsPythonCode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     "print(2 + 2)",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stdout != "4\n" {
		t.Fatalf("expected stdout %q, got %q", "4\n", result.Stdout)
	}

	if result.Stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.Stderr)
	}

	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", result.ExitCode)
	}
}
func TestDockerExecutorRunsJavaScriptCode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, Request{
		Language: "javascript",
		Code:     "console.log(2 + 2)",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stdout != "4\n" {
		t.Fatalf("expected stdout %q, got %q", "4\n", result.Stdout)
	}

	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", result.ExitCode)
	}
}

func TestDockerExecutorReturnsTimeoutError(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutor()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     "while True: pass",
	})

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
}

func TestDockerExecutorAppliesOutputLimitFromConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}

	executor := NewDockerExecutorWithConfig(DockerConfig{
		Memory:     "128m",
		CPUs:       "0.5",
		OutputSize: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := executor.Run(ctx, Request{
		Language: "python",
		Code:     `print("x" * 100)`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Stdout) != 10 {
		t.Fatalf("expected stdout length 10, got %d", len(result.Stdout))
	}
}
