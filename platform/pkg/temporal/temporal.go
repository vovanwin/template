// Package temporal предоставляет удобную обертку для работы с Temporal.io
package temporal

import (
	"context"
	"fmt"

	"github.com/vovanwin/platform/pkg/logger"
	"go.uber.org/zap"
)

// Service основной сервис для работы с Temporal
type Service struct {
	client *Client
	worker *Worker
}

// ServiceConfig полная конфигурация Temporal сервиса
type ServiceConfig struct {
	Client Config       `yaml:"client"`
	Worker WorkerConfig `yaml:"worker"`
}

// NewService создает новый Temporal сервис
func NewService(config ServiceConfig) (*Service, error) {
	lg := logger.Named("temporal-service")

	// Создаем клиент
	client, err := NewClient(config.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	// Создаем воркер
	worker := NewWorker(client, config.Worker)

	service := &Service{
		client: client,
		worker: worker,
	}

	lg.Info(context.Background(), "Temporal service created",
		zap.String("host", config.Client.Host),
		zap.Int("port", config.Client.Port),
		zap.String("namespace", config.Client.Namespace),
		zap.String("task_queue", config.Worker.TaskQueue))

	return service, nil
}

// GetClient возвращает клиент
func (s *Service) GetClient() *Client {
	return s.client
}

// GetWorker возвращает воркер
func (s *Service) GetWorker() *Worker {
	return s.worker
}

// Start запускает сервис
func (s *Service) Start(ctx context.Context) error {
	logger.Info(ctx, "Starting Temporal service")

	// Запускаем воркер
	if err := s.worker.Start(ctx); err != nil {
		return fmt.Errorf("failed to start temporal worker: %w", err)
	}

	return nil
}

// Stop останавливает сервис
func (s *Service) Stop(ctx context.Context) {
	logger.Info(ctx, "Stopping Temporal service")

	s.worker.Stop(ctx)
	s.client.Close()
}
