package workflows

import (
	"github.com/vovanwin/template/app/internal/workflows/activities"
	"github.com/vovanwin/template/app/internal/workflows/workflows"

	"github.com/vovanwin/platform/pkg/temporal"
	"go.uber.org/fx"
)

// Module модуль для регистрации воркфлоу и активностей
var Module = fx.Module("workflows",
	fx.Provide(
		activities.NewUserActivities,
	),
	fx.Invoke(RegisterWorkflows),
)

// RegisterWorkflows регистрирует все воркфлоу и активности в Temporal воркере
func RegisterWorkflows(temporalService *temporal.Service, userActivities *activities.UserActivities) {
	worker := temporalService.GetWorker()

	// Регистрируем воркфлоу
	worker.RegisterWorkflow(workflows.UserOnboardingWorkflow)

	// Регистрируем активности с именами для удобства вызова
	worker.RegisterActivityWithName(userActivities.ValidateUserDataActivity, "ValidateUserDataActivity")
	worker.RegisterActivityWithName(userActivities.CreateUserProfileActivity, "CreateUserProfileActivity")
	worker.RegisterActivityWithName(userActivities.SendWelcomeEmailActivity, "SendWelcomeEmailActivity")
	worker.RegisterActivityWithName(userActivities.SendNotificationActivity, "SendNotificationActivity")
}
