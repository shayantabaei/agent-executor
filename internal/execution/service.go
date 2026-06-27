package execution

import (
	"context"
	"time"
)

type ServiceConfig struct {
	Timeout time.Duration
}

func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		Timeout: 5 * time.Second,
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
	if config.Timeout == 0 {
		config.Timeout = DefaultServiceConfig().Timeout
	}

	return &Service{
		executor: executor,
		config:   config,
	}
}

func (s *Service) Run(ctx context.Context, req Request) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	return s.executor.Run(ctx, req)
}
