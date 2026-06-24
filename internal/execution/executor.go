package execution

import "context"

type Request struct {
	Language string
	Code     string
	Files    []InputFile
}

type InputFile struct {
	Path    string
	Content string
}

type Result struct {
	Stdout    string
	Stderr    string
	ExitCode  int
	Artifacts []Artifact
}

type Artifact struct {
	Path        string
	Size        int64
	Content     string
	Encoding    string
	ContentType string
}

// Executor executes code in a supported programming language.
type Executor interface {
	Run(ctx context.Context, req Request) (Result, error)
}
