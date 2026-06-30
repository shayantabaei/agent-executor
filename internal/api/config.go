package api

type Config struct {
	MaxBodySize int64
}

func DefaultConfig() Config {
	return Config{
		// Allow a small amount of JSON overhead around the code payload (80 KiB)
		MaxBodySize: 80 * 1024,
	}
}
