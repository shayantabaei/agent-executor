package api

import "testing"

func TestToExecutionRequestMapsFields(t *testing.T) {
	req := ExecutionRequest{
		Language: "python",
		Code:     "print('hello')",
		Files: []InputFile{
			{
				Path:    "data/input.txt",
				Content: "hello",
			},
		},
	}

	executionReq := toExecutionRequest(req)

	if executionReq.Language != req.Language {
		t.Fatalf("expected language %q, got %q", req.Language, executionReq.Language)
	}

	if executionReq.Code != req.Code {
		t.Fatalf("expected code %q, got %q", req.Code, executionReq.Code)
	}

	if len(executionReq.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(executionReq.Files))
	}

	if executionReq.Files[0].Path != "data/input.txt" {
		t.Fatalf("expected file path %q, got %q", "data/input.txt", executionReq.Files[0].Path)
	}

	if executionReq.Files[0].Content != "hello" {
		t.Fatalf("expected file content %q, got %q", "hello", executionReq.Files[0].Content)
	}
}
