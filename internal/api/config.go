package api

type Config struct {
	MaxBodySize      int64
	MaxCodeSize      int
	MaxFileCount     int
	MaxFileSizeBytes int
	MaxTotalFileSize int
}

func DefaultConfig() Config {
	return Config{
		// Allow a small amount of JSON overhead around the code payload (80 KiB)
		MaxBodySize: 80 * 1024,
		// Reject code payloads larger than 64 KiB before execution.
		MaxCodeSize: 64 * 1024,
		// Limit the number of input files accepted in one execution request.
		MaxFileCount: 10,
		// Limit each individual input file to 32 KiB.
		MaxFileSizeBytes: 32 * 1024,
		// Limit the combined size of all input files to 64 KiB.
		MaxTotalFileSize: 64 * 1024,
	}
}
