package api

type ExecutionRequest struct {
	Language string      `json:"language"`
	Code     string      `json:"code"`
	Files    []InputFile `json:"files,omitempty"`
}

type InputFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type ExecutionResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

type RuntimesResponse struct {
	Runtimes []string `json:"runtimes"`
}
