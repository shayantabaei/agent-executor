package execution

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeExecutor struct {
	result Result
	err    error
	called bool
}

func (f *fakeExecutor) Run(ctx context.Context, req Request) (Result, error) {
	f.called = true
	return f.result, f.err
}

func TestServiceRunRejectsMissingLanguage(t *testing.T) {
	executor := &fakeExecutor{}

	service := NewServiceWithConfig(executor, ServiceConfig{
		Timeout:          time.Second,
		MaxCodeSize:      64 * 1024,
		MaxFileCount:     10,
		MaxFileSizeBytes: 64 * 1024,
		MaxTotalFileSize: 256 * 1024,
	})

	_, err := service.Run(context.Background(), Request{
		Code: "print(2 + 2)",
	})

	if !errors.Is(err, ErrLanguageRequired) {
		t.Fatalf("expected ErrLanguageRequired, got %v", err)
	}

	if executor.called {
		t.Fatal("expected executor not to be called")
	}
}

func TestServiceRunCallsExecutorForValidRequest(t *testing.T) {
	executor := &fakeExecutor{
		result: Result{
			Stdout:   "ok\n",
			ExitCode: 0,
		},
	}

	service := NewServiceWithConfig(executor, ServiceConfig{
		Timeout:          time.Second,
		MaxCodeSize:      64 * 1024,
		MaxFileCount:     10,
		MaxFileSizeBytes: 64 * 1024,
		MaxTotalFileSize: 256 * 1024,
	})

	result, err := service.Run(context.Background(), Request{
		Language: "python",
		Code:     "print('ok')",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !executor.called {
		t.Fatal("expected executor to be called")
	}

	if result.Stdout != "ok\n" {
		t.Fatalf("expected stdout %q, got %q", "ok\n", result.Stdout)
	}
}

type blockingExecutor struct{}

func (blockingExecutor) Run(ctx context.Context, req Request) (Result, error) {
	<-ctx.Done()
	return Result{}, ctx.Err()
}

func TestServiceRunAppliesTimeout(t *testing.T) {
	service := NewServiceWithConfig(blockingExecutor{}, ServiceConfig{
		Timeout:          10 * time.Millisecond,
		MaxCodeSize:      64 * 1024,
		MaxFileCount:     10,
		MaxFileSizeBytes: 64 * 1024,
		MaxTotalFileSize: 256 * 1024,
	})

	_, err := service.Run(context.Background(), Request{
		Language: "python",
		Code:     "print('hello')",
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
}
