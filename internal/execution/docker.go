package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
)

type DockerExecutor struct{}

func NewDockerExecutor() *DockerExecutor {
	return &DockerExecutor{}
}

func (e *DockerExecutor) Run(
	ctx context.Context,
	req Request,
) (Result, error) {

	runtime, err := runtimeForLanguage(req.Language)

	if err != nil {
		return Result{}, err
	}

	args := []string{
		"run",
		"--rm",
		"-i",
		"--network", "none",
		"--memory", "128m",
		"--cpus", "0.5",
		runtime.Image(),
	}

	args = append(args, runtime.Command()...)

	// Create the command to run the Docker container with the specified arguments.
	cmd := exec.CommandContext(ctx, "docker", args...)
	// Provide the code to the container via stdin.
	cmd.Stdin = bytes.NewBufferString(req.Code)

	// Set a limit of 64KB
	stdout := NewLimitedWriter(64 * 1024)
	stderr := NewLimitedWriter(64 * 1024)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Run the docker command
	err = cmd.Run()

	result := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if ctx.Err() != nil {
		return Result{}, ctx.Err()
	}

	if err == nil {
		return result, nil
	}

	var exitError *exec.ExitError
	if ok := errors.As(err, &exitError); ok {
		result.ExitCode = exitError.ExitCode()
		return result, nil
	}

	return Result{}, fmt.Errorf("failed to execute code: %w", err)

}
