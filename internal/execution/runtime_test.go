package execution

import "testing"

func TestRuntimeForLanguageReturnsPythonRuntime(t *testing.T) {
	runtime, err := runtimeForLanguage("python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if runtime.Image() != "python:3.12-alpine" {
		t.Fatalf("expected python image, got %q", runtime.Image())
	}
}

func TestRuntimeForLanguageReturnsUnsupportedLanguageError(t *testing.T) {
	_, err := runtimeForLanguage("ruby")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(UnsupportedLanguageError); !ok {
		t.Fatalf("expected UnsupportedLanguageError, got %T", err)
	}
}
