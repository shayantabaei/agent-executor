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
	Stdout    string             `json:"stdout"`
	Stderr    string             `json:"stderr"`
	ExitCode  int                `json:"exitCode"`
	Artifacts []ArtifactResponse `json:"artifacts,omitempty"`
}

type ArtifactResponse struct {
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	Content     string `json:"content,omitempty"`
	Encoding    string `json:"encoding,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

type RuntimesResponse struct {
	Runtimes []string `json:"runtimes"`
}
