package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

type DockerExecutor struct {
	config DockerConfig
}

func NewDockerExecutor() *DockerExecutor {
	return NewDockerExecutorWithConfig(DefaultDockerConfig())
}

func NewDockerExecutorWithConfig(config DockerConfig) *DockerExecutor {
	return &DockerExecutor{config: config}
}

func (e *DockerExecutor) Run(
	ctx context.Context,
	req Request,
) (Result, error) {

	runtime, err := runtimeForLanguage(req.Language)

	if err != nil {
		return Result{}, err
	}

	ws, err := createWorkspace(req.Files)
	if err != nil {
		return Result{}, err
	}
	defer ws.cleanup()

	args := e.buildRunArgs(runtime, ws.path)

	// Create the command to run the Docker container with the specified arguments.
	cmd := exec.CommandContext(ctx, "docker", args...)
	// Provide the code to the container via stdin.
	cmd.Stdin = bytes.NewBufferString(req.Code)

	// Use config limit
	stdout := NewLimitedWriter(e.config.OutputSize)
	stderr := NewLimitedWriter(e.config.OutputSize)

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

func (e *DockerExecutor) buildRunArgs(runtime Runtime, workspacePath string) []string {
	cfg := e.config

	args := []string{
		"run",
		"--rm",
		"-i",
	}

	if cfg.NetworkDisabled {
		args = append(args, "--network", "none")
	}

	if cfg.Memory != "" {
		args = append(args, "--memory", cfg.Memory)
	}

	if cfg.CPUs != "" {
		args = append(args, "--cpus", cfg.CPUs)
	}

	if cfg.NoNewPrivileges {
		args = append(args, "--security-opt", "no-new-privileges:true")
	}

	if cfg.PidsLimit > 0 {
		args = append(args, "--pids-limit", strconv.Itoa(cfg.PidsLimit))
	}

	if cfg.PullPolicy != "" {
		args = append(args, "--pull", cfg.PullPolicy)
	}

	const containerWorkspacePath = "/workspace"

	if workspacePath != "" {
		// Mount the temporary host workspace into the container.
		// --mount avoids ambiguity with Windows drive-letter paths like C:\...
		args = append(
			args,
			"--mount",
			"type=bind,source="+workspacePath+",target="+containerWorkspacePath,
		)

		// Run the user code from /workspace so relative file paths work.
		args = append(args, "-w", containerWorkspacePath)
	}

	args = append(args, runtime.Image())
	args = append(args, runtime.Command()...)

	return args
}
