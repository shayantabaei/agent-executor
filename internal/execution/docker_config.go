package execution

type DockerConfig struct {
	Memory          string
	CPUs            string
	OutputSize      int
	NetworkDisabled bool
	NoNewPrivileges bool
	PidsLimit       int
	PullPolicy      string
}

func DefaultDockerConfig() DockerConfig {
	return DockerConfig{
		// Keep memory small so runaway programs cannot consume the host.
		Memory: "128m",

		// Limit CPU usage to half a core to reduce impact from busy loops.
		CPUs: "0.5",

		// Capture at most 64 KiB each for stdout and stderr.
		// This prevents unbounded output from growing memory usage.
		OutputSize: 64 * 1024,

		// Disable network access by default for local code execution.
		NetworkDisabled: true,

		// Prevent processes from gaining additional privileges inside the container.
		NoNewPrivileges: true,

		// Limit # of processors container can create (guard against fork bombs)
		PidsLimit: 128,

		// Prevent docker run from pulling images from a registry
		PullPolicy: "never",
	}
}
