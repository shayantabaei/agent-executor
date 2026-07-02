package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

type DockerExecutor struct {
	config DockerConfig

	createWorkspace  func(files []InputFile) (*workspace, error)
	cleanupWorkspace func(ws *workspace)
}

func NewDockerExecutor() *DockerExecutor {
	return NewDockerExecutorWithConfig(DefaultDockerConfig())
}

func NewDockerExecutorWithConfig(config DockerConfig) *DockerExecutor {
	return &DockerExecutor{
		config:           config,
		createWorkspace:  createWorkspace,
		cleanupWorkspace: cleanupWorkspace,
	}
}

func (e *DockerExecutor) Run(
	ctx context.Context,
	req Request,
) (Result, error) {

	runtime, err := runtimeForLanguage(req.Language)

	if err != nil {
		return Result{}, err
	}

	ws, err := e.createWorkspace(req.Files)
	if err != nil {
		return Result{}, err
	}
	defer e.cleanupWorkspace(ws)

	containerName := fmt.Sprintf("agent-executor-%d", time.Now().UnixNano())

	args := e.buildRunArgs(runtime, ws.path, containerName)

	// Create the command to run the Docker container with the specified arguments.
	cmd := exec.CommandContext(ctx, "docker", args...)
	defer func() {
		_ = exec.Command("docker", "rm", "-f", containerName).Run()
	}()
	// Provide the code to the container via stdin.
	cmd.Stdin = bytes.NewBufferString(req.Code)

	// Use config limit
	stdout := NewLimitedWriter(e.config.OutputSize)
	stderr := NewLimitedWriter(e.config.OutputSize)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Run the docker command
	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		_ = exec.Command("docker", "rm", "-f", containerName).Run()
		return Result{}, ctx.Err()
	}

	if ctx.Err() != nil {
		_ = exec.Command("docker", "rm", "-f", containerName).Run()
		return Result{}, ctx.Err()
	}

	artifacts, artifactError := ws.collectArtifacts(req.Files, e.config)

	if artifactError != nil {
		return Result{
			Stdout:    stdout.String(),
			Stderr:    stderr.String(),
			ErrorType: ErrorTypeArtifact,
		}, artifactError
	}

	result := Result{
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		Artifacts: artifacts,
	}

	if err == nil {
		return result, nil
	}

	var exitError *exec.ExitError
	if ok := errors.As(err, &exitError); ok {
		result.ExitCode = exitError.ExitCode()
		result.ErrorType = ErrorTypeRuntime
		return result, nil
	}

	return Result{}, fmt.Errorf("failed to execute code: %w", err)

}

func (e *DockerExecutor) buildRunArgs(runtime Runtime, workspacePath string, containerName string) []string {
	cfg := e.config

	args := []string{
		"run",
		"--rm",
		"-i",
	}

	if containerName != "" {
		args = append(args, "--name", containerName)
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
