package temporal

import (
	"context"
	"log/slog"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
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
func NewWorker(client *Client, config WorkerConfig, slogLogger *slog.Logger) *Worker {
	temporalWorker := worker.New(client.GetClient(), config.TaskQueue, worker.Options{
		Logger: log.NewStructuredLogger(slogLogger),
	})

	return &Worker{
		temporalWorker: temporalWorker,
		taskQueue:      config.TaskQueue,
	}
}

// RegisterWorkflow регистрирует воркфлоу
func (w *Worker) RegisterWorkflow(workflows ...interface{}) {
	for _, workflow := range workflows {
		w.temporalWorker.RegisterWorkflow(workflow)
	}
}

// RegisterActivity регистрирует активности
func (w *Worker) RegisterActivity(activities ...interface{}) {
	for _, activity := range activities {
		w.temporalWorker.RegisterActivity(activity)
	}
}

// RegisterWorkflowWithName регистрирует воркфлоу с именем
func (w *Worker) RegisterWorkflowWithName(workflow interface{}, name string) {
	// Просто используем стандартную регистрацию, Temporal сам определит имя
	w.temporalWorker.RegisterWorkflow(workflow)

}

// RegisterActivityWithName регистрирует активность с именем
func (w *Worker) RegisterActivityWithName(activity interface{}, name string) {
	// Просто используем стандартную регистрацию, Temporal сам определит имя
	w.temporalWorker.RegisterActivity(activity)

}

// GetRegistry возвращает реестр для регистрации воркфлоу и активностей
func (w *Worker) GetRegistry() worker.Worker {
	return w.temporalWorker
}

// Start запускает воркер
func (w *Worker) Start(ctx context.Context) error {

	// Запускаем воркер в горутине
	go func() {
		if err := w.temporalWorker.Run(worker.InterruptCh()); err != nil {

		}
	}()

	return nil
}

// Stop останавливает воркер
func (w *Worker) Stop(ctx context.Context) {
	w.temporalWorker.Stop()
}
