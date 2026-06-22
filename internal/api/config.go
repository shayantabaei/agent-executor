package api

type Config struct {
	MaxBodySize int64
	MaxCodeSize int
}

func DefaultConfig() Config {
	return Config{
		// Allow a small amount of JSON overhead around the code payload (80 KiB)
		MaxBodySize: 80 * 1024,
		// Reject code payloads larger than 64 KiB before execution.
		MaxCodeSize: 64 * 1024,
	}
}
