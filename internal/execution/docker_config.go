package execution

type DockerConfig struct {
	Memory                    string
	CPUs                      string
	OutputSize                int
	NetworkDisabled           bool
	NoNewPrivileges           bool
	PidsLimit                 int
	PullPolicy                string
	MaxArtifactCount          int
	MaxArtifactSizeBytes      int64
	MaxTotalArtifactSizeBytes int64
	MaxInlineArtifactBytes    int64
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

		// Limit the number of processes the container can create.
		// This helps reduce the impact of fork bombs.
		PidsLimit: 128,

		// Prevent docker run from implicitly pulling images from a registry.
		// Runtime images must already exist locally.
		PullPolicy: "never",

		// Limit the number of generated files returned as artifacts.
		MaxArtifactCount: 10,

		// Reject any single artifact larger than 64 KiB.
		MaxArtifactSizeBytes: 64 * 1024,

		// Reject artifact collection if all artifacts exceed 256 KiB combined.
		MaxTotalArtifactSizeBytes: 256 * 1024,

		// Inline content only for artifacts up to 32 KiB.
		// Larger artifacts can still be returned as metadata-only results.
		MaxInlineArtifactBytes: 32 * 1024,
	}
}
