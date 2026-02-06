package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

// Client обертка над temporal клиентом с дополнительной функциональностью
type Client struct {
	temporalClient client.Client
	namespace      string
}

// Config конфигурация для подключения к Temporal
type Config struct {
	Host      string `yaml:"host" env:"TEMPORAL_HOST" validate:"required"`
	Port      int    `yaml:"port" env:"TEMPORAL_PORT" default:"7233"`
	Namespace string `yaml:"namespace" env:"TEMPORAL_NAMESPACE" default:"default"`
}

// NewClient создает новый Temporal клиент
func NewClient(config Config) (*Client, error) {
	clientOptions := client.Options{
		HostPort:  fmt.Sprintf("%s:%d", config.Host, config.Port),
		Namespace: config.Namespace,
	}

	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to temporal: %w", err)
	}

	return &Client{
		temporalClient: temporalClient,
		namespace:      config.Namespace,
	}, nil
}

// GetClient возвращает базовый temporal клиент
func (c *Client) GetClient() client.Client {
	return c.temporalClient
}

// GetNamespace возвращает namespace
func (c *Client) GetNamespace() string {
	return c.namespace
}

// Close закрывает соединение
func (c *Client) Close() {
	c.temporalClient.Close()
}

// ExecuteWorkflow запускает воркфлоу
func (c *Client) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return c.temporalClient.ExecuteWorkflow(ctx, options, workflow, args...)
}

// GetWorkflow получает информацию о воркфлоу
func (c *Client) GetWorkflow(ctx context.Context, workflowID, runID string) client.WorkflowRun {
	return c.temporalClient.GetWorkflow(ctx, workflowID, runID)
}

// ScheduleWorkflow планирует выполнение воркфлоу
func (c *Client) ScheduleWorkflow(ctx context.Context, scheduleID string, schedule client.ScheduleSpec, action *client.ScheduleWorkflowAction, opts client.ScheduleOptions) (client.ScheduleHandle, error) {
	return c.temporalClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID:   scheduleID,
		Spec: schedule,
		Action: &client.ScheduleWorkflowAction{
			ID:                       action.ID,
			Workflow:                 action.Workflow,
			Args:                     action.Args,
			TaskQueue:                action.TaskQueue,
			WorkflowExecutionTimeout: action.WorkflowExecutionTimeout,
			WorkflowTaskTimeout:      action.WorkflowTaskTimeout,
		},
	})
}
