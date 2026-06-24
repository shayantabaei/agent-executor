package api

import (
	"errors"
	"testing"
)

func TestValidateExecutionRequestAllowsValidRequestWithoutFiles(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
	}

	if err := validateExecutionRequest(req, cfg); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestValidateExecutionRequestAllowsSafeNestedFilePaths(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{
				Path:    "src/main.py",
				Content: "print('hello')",
			},
			{
				Path:    "data/input.txt",
				Content: "hello",
			},
		},
	}

	if err := validateExecutionRequest(req, cfg); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestValidateExecutionRequestRejectsMissingLanguage(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Code: "print('hello')",
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrLanguageRequired) {
		t.Fatalf("expected ErrLanguageRequired, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsMissingCode(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrCodeRequired) {
		t.Fatalf("expected ErrCodeRequired, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsCodeOverLimit(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxCodeSize = 3

	req := ExecutionRequest{
		Language: "python",
		Code:     "abcd",
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrCodeTooLarge) {
		t.Fatalf("expected ErrCodeTooLarge, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsTooManyFiles(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxFileCount = 1

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "a.txt", Content: "a"},
			{Path: "b.txt", Content: "b"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrTooManyFiles) {
		t.Fatalf("expected ErrTooManyFiles, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsLargeFile(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxFileSizeBytes = 3

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "a.txt", Content: "abcd"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected ErrFileTooLarge, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsLargeTotalFileSize(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxFileSizeBytes = 10
	cfg.MaxTotalFileSize = 5

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "a.txt", Content: "abc"},
			{Path: "b.txt", Content: "def"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrTotalFilesTooLarge) {
		t.Fatalf("expected ErrTotalFilesTooLarge, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsAbsoluteFilePath(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "/etc/passwd", Content: "bad"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrInvalidFilePath) {
		t.Fatalf("expected ErrInvalidFilePath, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsPathTraversal(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "../secret.txt", Content: "bad"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrInvalidFilePath) {
		t.Fatalf("expected ErrInvalidFilePath, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsNestedPathTraversal(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "safe/../../secret.txt", Content: "bad"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrInvalidFilePath) {
		t.Fatalf("expected ErrInvalidFilePath, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsEmptyFilePath(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: "", Content: "bad"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrInvalidFilePath) {
		t.Fatalf("expected ErrInvalidFilePath, got %v", err)
	}
}

func TestValidateExecutionRequestRejectsBackslashFilePath(t *testing.T) {
	cfg := DefaultConfig()

	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{Path: `data\input.txt`, Content: "bad"},
		},
	}

	err := validateExecutionRequest(req, cfg)
	if !errors.Is(err, ErrInvalidFilePath) {
		t.Fatalf("expected ErrInvalidFilePath, got %v", err)
	}
}
