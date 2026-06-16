package execution

import "context"

type Request struct {
	Language string
	Code     string
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Executor executes code in a supported programming language.
type Executor interface {
	Run(ctx context.Context, req Request) (Result, error)
}
