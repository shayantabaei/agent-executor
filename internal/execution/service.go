package execution

import (
	"context"
	"errors"
	"time"
)

type ServiceConfig struct {
	Timeout          time.Duration
	MaxCodeSize      int
	MaxFileCount     int
	MaxFileSizeBytes int
	MaxTotalFileSize int
}

func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		Timeout: 5 * time.Second,
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

type Service struct {
	executor Executor
	config   ServiceConfig
}

func NewService(executor Executor) *Service {
	return NewServiceWithConfig(executor, DefaultServiceConfig())
}

func NewServiceWithConfig(executor Executor, config ServiceConfig) *Service {

	defaults := DefaultServiceConfig()

	if config.Timeout == 0 {
		config.Timeout = defaults.Timeout
	}

	if config.MaxCodeSize == 0 {
		config.MaxCodeSize = defaults.MaxCodeSize
	}

	if config.MaxFileCount == 0 {
		config.MaxFileCount = defaults.MaxFileCount
	}

	if config.MaxFileSizeBytes == 0 {
		config.MaxFileSizeBytes = defaults.MaxFileSizeBytes
	}

	if config.MaxTotalFileSize == 0 {
		config.MaxTotalFileSize = defaults.MaxTotalFileSize
	}

	return &Service{
		executor: executor,
		config:   config,
	}
}

func (s *Service) Run(ctx context.Context, req Request) (Result, error) {
	start := time.Now()

	if err := validateRequest(req, s.config); err != nil {
		return Result{
			DurationMs: time.Since(start).Milliseconds(),
			ErrorType:  ErrorTypeValidation,
		}, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	result, err := s.executor.Run(ctx, req)
	result.DurationMs = time.Since(start).Milliseconds()

	if err == nil {
		return result, nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		result.TimedOut = true
		if result.Stderr == "" {
			result.Stderr = "execution timed out"
		}
		result.ErrorType = ErrorTypeTimeout
		return result, err
	}

	if result.ErrorType == "" {
		result.ErrorType = classifyError(err)
	}

	return result, err
}

func classifyError(err error) ErrorType {
	var unsupportedLanguageError UnsupportedLanguageError
	if errors.As(err, &unsupportedLanguageError) {
		return ErrorTypeRuntimeNotFound
	}

	if IsValidationError(err) {
		return ErrorTypeValidation
	}

	return ErrorTypeInternal
}
