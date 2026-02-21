package workflows

import (
	"github.com/vovanwin/template/internal/pkg/temporal"
	"github.com/vovanwin/template/internal/workflows/activities"
	"github.com/vovanwin/template/internal/workflows/reminder"
	"github.com/vovanwin/template/internal/workflows/workflows"
	reminderv1 "github.com/vovanwin/template/pkg/temporal/reminder"

	"go.uber.org/fx"
)

// Module модуль для регистрации воркфлоу и активностей
var Module = fx.Module("workflows",
	fx.Provide(
		activities.NewUserActivities,
		reminder.NewWorkflows,
		reminder.NewActivities,
	),
	fx.Invoke(RegisterWorkflows),
)

// RegisterWorkflows регистрирует все воркфлоу и активности в Temporal воркере
func RegisterWorkflows(
	temporalService *temporal.Service,
	userActivities *activities.UserActivities,
	reminderWorkflows *reminder.Workflows,
	reminderActivities *reminder.Activities,
) {
	worker := temporalService.GetWorker()

	// User onboarding
	worker.RegisterWorkflow(workflows.UserOnboardingWorkflow)
	worker.RegisterActivityWithName(userActivities.ValidateUserDataActivity, "ValidateUserDataActivity")
	worker.RegisterActivityWithName(userActivities.CreateUserProfileActivity, "CreateUserProfileActivity")
	worker.RegisterActivityWithName(userActivities.SendWelcomeEmailActivity, "SendWelcomeEmailActivity")
	worker.RegisterActivityWithName(userActivities.SendNotificationActivity, "SendNotificationActivity")

	// Reminder workflows
	reminderv1.RegisterReminderWorkflows(worker.GetRegistry(), reminderWorkflows)
	reminderv1.RegisterReminderActivities(worker.GetRegistry(), reminderActivities)
}
