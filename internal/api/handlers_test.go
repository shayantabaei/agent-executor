package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shayantabaei/agent-executor/internal/execution"
)

type stubExecutor struct {
	result execution.Result
	err    error
}

func (s stubExecutor) Run(
	_ context.Context,
	_ execution.Request,
) (execution.Result, error) {
	return s.result, s.err
}

func TestExecutionHandlerSuccess(t *testing.T) {
	executor := stubExecutor{
		result: execution.Result{
			Stdout:   "4\n",
			Stderr:   "",
			ExitCode: 0,
		},
		err: nil,
	}

	handler := NewHandler(executor)

	requestBody := `{"language": "python", "code": "print(2 + 2)"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/executions",
		strings.NewReader(requestBody),
	)

	recorder := httptest.NewRecorder()
	handler.ExecutionHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	expected := `{"stdout":"4\n","stderr":"","exitCode":0}`

	if strings.TrimSpace(recorder.Body.String()) != expected {
		t.Fatalf(
			"expected body %q, got %q",
			expected,
			recorder.Body.String(),
		)
	}

}

func TestExecutionHandlerRejectsMissingLanguage(t *testing.T) {
	handler := NewHandler(stubExecutor{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/executions",
		strings.NewReader(`{"code":"print(2 + 2)"}`),
	)

	recorder := httptest.NewRecorder()

	handler.ExecutionHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}

func TestExecutionHandlerReturnsInternalServerError(t *testing.T) {
	handler := NewHandler(stubExecutor{
		err: errors.New("docker unavailable"),
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/executions",
		strings.NewReader(`{
			"language":"python",
			"code":"print(2 + 2)"
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.ExecutionHandler(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}

func TestExecutionHandlerRejectsCodeOverLimit(t *testing.T) {
	handler := NewHandlerWithConfig(stubExecutor{}, Config{
		MaxBodySize: 1024,
		MaxCodeSize: 10,
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/executions",
		strings.NewReader(`{
			"language": "python",
			"code": "this code is definitely too long"
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.ExecutionHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestExecutionHandlerRejectsBodyOverLimit(t *testing.T) {
	handler := NewHandlerWithConfig(stubExecutor{}, Config{
		MaxBodySize: 10,
		MaxCodeSize: 64 * 1024,
	})

	req := httptest.NewRequest(
		http.MethodPost,
		"/executions",
		strings.NewReader(`{
			"language": "python",
			"code": "print(2 + 2)"
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.ExecutionHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
