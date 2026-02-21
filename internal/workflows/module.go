package workflows

import (
	"github.com/vovanwin/template/internal/pkg/temporal"
	"github.com/vovanwin/template/internal/workflows/reminder"
	reminderv1 "github.com/vovanwin/template/pkg/temporal/reminder"

	"go.uber.org/fx"
)

// Module модуль для регистрации воркфлоу и активностей
var Module = fx.Module("workflows",
	fx.Provide(
		reminder.NewWorkflows,
		reminder.NewActivities,
	),
	fx.Invoke(RegisterWorkflows),
)

// RegisterWorkflows регистрирует все воркфлоу и активности в Temporal воркере
func RegisterWorkflows(
	temporalService *temporal.Service,
	reminderWorkflows *reminder.Workflows,
	reminderActivities *reminder.Activities,
) {
	worker := temporalService.GetWorker()

	// Reminder workflows
	reminderv1.RegisterReminderWorkflows(worker.GetRegistry(), reminderWorkflows)
	reminderv1.RegisterReminderActivities(worker.GetRegistry(), reminderActivities)
}
