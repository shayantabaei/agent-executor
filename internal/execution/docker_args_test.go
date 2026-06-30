package execution

import "testing"

func TestDockerExecutorBuildRunArgsUsesConfig(t *testing.T) {
	executor := NewDockerExecutorWithConfig(DockerConfig{
		Memory:          "128m",
		CPUs:            "0.5",
		OutputSize:      64 * 1024,
		NetworkDisabled: true,
		NoNewPrivileges: true,
		PidsLimit:       128,
		PullPolicy:      "never",
	})

	runtime, err := runtimeForLanguage("python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	args := executor.buildRunArgs(runtime, "/tmp/workspace", "test-container")

	expectedPrefix := []string{
		"run",
		"--rm",
		"-i",
		"--name", "test-container",
		"--network", "none",
		"--memory", "128m",
		"--cpus", "0.5",
		"--security-opt", "no-new-privileges:true",
		"--pids-limit", "128",
		"--pull", "never",
		"--mount", "type=bind,source=/tmp/workspace,target=/workspace",
		"-w", "/workspace",
		runtime.Image(),
	}

	if len(args) < len(expectedPrefix) {
		t.Fatalf("expected args to have at least %d entries, got %d: %v", len(expectedPrefix), len(args), args)
	}

	for i, expectedArg := range expectedPrefix {
		if args[i] != expectedArg {
			t.Fatalf("expected arg %q at index %d, got %q. args: %v", expectedArg, i, args[i], args)
		}
	}
}

func TestDefaultDockerConfigIncludesHardeningDefaults(t *testing.T) {
	cfg := DefaultDockerConfig()

	if cfg.Memory != "128m" {
		t.Fatalf("expected memory %q, got %q", "128m", cfg.Memory)
	}

	if cfg.CPUs != "0.5" {
		t.Fatalf("expected CPUs %q, got %q", "0.5", cfg.CPUs)
	}

	if cfg.OutputSize != 64*1024 {
		t.Fatalf("expected output size %d, got %d", 64*1024, cfg.OutputSize)
	}

	if !cfg.NetworkDisabled {
		t.Fatal("expected network to be disabled by default")
	}

	if !cfg.NoNewPrivileges {
		t.Fatal("expected no-new-privileges to be enabled by default")
	}

	if cfg.PidsLimit != 128 {
		t.Fatalf("expected PID limit %d, got %d", 128, cfg.PidsLimit)
	}

	if cfg.PullPolicy != "never" {
		t.Fatalf("expected pull policy %q, got %q", "never", cfg.PullPolicy)
	}
}
