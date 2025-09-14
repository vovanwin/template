package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// UserOnboardingWorkflowInput входные данные для воркфлоу регистрации пользователя
type UserOnboardingWorkflowInput struct {
	UserID   string                 `json:"user_id"`
	Email    string                 `json:"email"`
	UserData map[string]interface{} `json:"user_data"`
}

// UserOnboardingWorkflowResult результат выполнения воркфлоу регистрации пользователя
type UserOnboardingWorkflowResult struct {
	Success   bool   `json:"success"`
	UserID    string `json:"user_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// UserOnboardingWorkflow воркфлоу для регистрации и настройки нового пользователя
func UserOnboardingWorkflow(ctx workflow.Context, input UserOnboardingWorkflowInput) (*UserOnboardingWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting user onboarding workflow", "user_id", input.UserID)

	// Настройки активностей
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Шаг 1: Валидация данных пользователя
	logger.Info("Step 1: Validating user data", "user_id", input.UserID)
	err := workflow.ExecuteActivity(ctx, "ValidateUserDataActivity", input.UserData).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to validate user data", "user_id", input.UserID, "error", err)
		return &UserOnboardingWorkflowResult{
			Success:   false,
			UserID:    input.UserID,
			Message:   "User data validation failed: " + err.Error(),
			Timestamp: workflow.Now(ctx).Format(time.RFC3339),
		}, err
	}

	// Шаг 2: Создание профиля пользователя
	logger.Info("Step 2: Creating user profile", "user_id", input.UserID)
	err = workflow.ExecuteActivity(ctx, "CreateUserProfileActivity", input.UserID, input.UserData).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to create user profile", "user_id", input.UserID, "error", err)
		return &UserOnboardingWorkflowResult{
			Success:   false,
			UserID:    input.UserID,
			Message:   "Profile creation failed: " + err.Error(),
			Timestamp: workflow.Now(ctx).Format(time.RFC3339),
		}, err
	}

	// Шаг 3: Отправка приветственного письма
	logger.Info("Step 3: Sending welcome email", "user_id", input.UserID)
	err = workflow.ExecuteActivity(ctx, "SendWelcomeEmailActivity", input.UserID, input.Email).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send welcome email", "user_id", input.UserID, "error", err)
		// Письмо не критично, продолжаем
		logger.Warn("Welcome email failed, continuing workflow", "user_id", input.UserID)
	}

	// Шаг 4: Отправка уведомления об успешной регистрации
	logger.Info("Step 4: Sending registration notification", "user_id", input.UserID)
	notificationMessage := "Welcome to our platform! Your account has been successfully created."
	err = workflow.ExecuteActivity(ctx, "SendNotificationActivity", input.UserID, notificationMessage).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send notification", "user_id", input.UserID, "error", err)
		// Уведомление не критично, продолжаем
		logger.Warn("Notification failed, continuing workflow", "user_id", input.UserID)
	}

	// Шаг 5: Задержка и отправка финального уведомления
	logger.Info("Step 5: Waiting before final notification", "user_id", input.UserID)
	workflow.Sleep(ctx, time.Second*30) // Ждем 30 секунд

	finalMessage := "Don't forget to complete your profile setup!"
	err = workflow.ExecuteActivity(ctx, "SendNotificationActivity", input.UserID, finalMessage).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send final notification", "user_id", input.UserID, "error", err)
	}

	logger.Info("User onboarding workflow completed successfully", "user_id", input.UserID)

	return &UserOnboardingWorkflowResult{
		Success:   true,
		UserID:    input.UserID,
		Message:   "User onboarding completed successfully",
		Timestamp: workflow.Now(ctx).Format(time.RFC3339),
	}, nil
}
