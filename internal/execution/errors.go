package execution

type ErrorType string

const (
	ErrorTypeValidation      ErrorType = "validation_error"
	ErrorTypeRuntimeNotFound ErrorType = "runtime_not_found"
	ErrorTypeDocker          ErrorType = "docker_error"
	ErrorTypeTimeout         ErrorType = "timeout"
	ErrorTypeWorkspace       ErrorType = "workspace_error"
	ErrorTypeArtifact        ErrorType = "artifact_error"
	ErrorTypeRuntime         ErrorType = "runtime_error"
	ErrorTypeInternal        ErrorType = "internal_error"
)
