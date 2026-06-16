package api

type ExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ExecutionResponse struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}
