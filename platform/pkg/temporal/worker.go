package temporal

import (
	"context"
	"fmt"

	"github.com/vovanwin/platform/pkg/logger"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

// Worker обертка над temporal воркером
type Worker struct {
	temporalWorker worker.Worker
	taskQueue      string
}

// WorkerConfig конфигурация воркера
type WorkerConfig struct {
	TaskQueue string `yaml:"task_queue" env:"TEMPORAL_TASK_QUEUE" validate:"required"`
}

// NewWorker создает новый Temporal воркер
func NewWorker(client *Client, config WorkerConfig) *Worker {
	temporalWorker := worker.New(client.GetClient(), config.TaskQueue, worker.Options{})

	return &Worker{
		temporalWorker: temporalWorker,
		taskQueue:      config.TaskQueue,
	}
}

// RegisterWorkflow регистрирует воркфлоу
func (w *Worker) RegisterWorkflow(workflows ...interface{}) {
	for _, workflow := range workflows {
		w.temporalWorker.RegisterWorkflow(workflow)
		logger.Info(context.Background(), "Registered workflow", zap.String("workflow", fmt.Sprintf("%T", workflow)))
	}
}

// RegisterActivity регистрирует активности
func (w *Worker) RegisterActivity(activities ...interface{}) {
	for _, activity := range activities {
		w.temporalWorker.RegisterActivity(activity)
		logger.Info(context.Background(), "Registered activity", zap.String("activity", fmt.Sprintf("%T", activity)))
	}
}

// RegisterWorkflowWithName регистрирует воркфлоу с именем
func (w *Worker) RegisterWorkflowWithName(workflow interface{}, name string) {
	// Просто используем стандартную регистрацию, Temporal сам определит имя
	w.temporalWorker.RegisterWorkflow(workflow)
	logger.Info(context.Background(), "Registered workflow with name",
		zap.String("workflow", fmt.Sprintf("%T", workflow)),
		zap.String("name", name))
}

// RegisterActivityWithName регистрирует активность с именем
func (w *Worker) RegisterActivityWithName(activity interface{}, name string) {
	// Просто используем стандартную регистрацию, Temporal сам определит имя
	w.temporalWorker.RegisterActivity(activity)
	logger.Info(context.Background(), "Registered activity with name",
		zap.String("activity", fmt.Sprintf("%T", activity)),
		zap.String("name", name))
}

// Start запускает воркер
func (w *Worker) Start(ctx context.Context) error {
	logger.Info(ctx, "Starting temporal worker", zap.String("task_queue", w.taskQueue))

	// Запускаем воркер в горутине
	go func() {
		if err := w.temporalWorker.Run(worker.InterruptCh()); err != nil {
			logger.Error(context.Background(), "Temporal worker error", zap.Error(err))
		}
	}()

	return nil
}

// Stop останавливает воркер
func (w *Worker) Stop(ctx context.Context) {
	logger.Info(ctx, "Stopping temporal worker", zap.String("task_queue", w.taskQueue))
	w.temporalWorker.Stop()
}
