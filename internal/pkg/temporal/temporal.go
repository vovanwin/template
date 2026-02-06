// Package temporal предоставляет удобную обертку для работы с Temporal.io
package temporal

import (
	"context"
	"fmt"
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

	// Запускаем воркер
	if err := s.worker.Start(ctx); err != nil {
		return fmt.Errorf("failed to start temporal worker: %w", err)
	}

	return nil
}

// Stop останавливает сервис
func (s *Service) Stop(ctx context.Context) {

	s.worker.Stop(ctx)
	s.client.Close()
}
