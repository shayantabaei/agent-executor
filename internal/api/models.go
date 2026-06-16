package api

type ExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ExecutionResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}
