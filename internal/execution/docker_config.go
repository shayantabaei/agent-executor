package execution

type DockerConfig struct {
	Memory     string
	CPUs       string
	OutputSize int
}

func DefaultDockerConfig() DockerConfig {
	return DockerConfig{
		Memory:     "128m",
		CPUs:       "0.5",
		OutputSize: 64 * 1024,
	}
}
