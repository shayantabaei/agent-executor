package execution

type DockerConfig struct {
	Memory     string
	CPUs       string
	OutputSize int
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
	}
}
