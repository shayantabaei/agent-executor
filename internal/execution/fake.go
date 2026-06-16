package execution

import "context"

// FakeExecutor is a temporary Executor implementation used before
// Docker-based execution is available.
type FakeExecutor struct{}

func (FakeExecutor) Run(
	_ context.Context,
	_ Request,
) (Result, error) {
	return Result{
		Stdout:   "execution not yet implemented\n",
		Stderr:   "",
		ExitCode: 0,
	}, nil
}
