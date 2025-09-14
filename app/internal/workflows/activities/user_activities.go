package activities

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
)

// UserActivities содержит активности для работы с пользователями
type UserActivities struct{}

// NewUserActivities создает новый экземпляр активностей пользователей
func NewUserActivities() *UserActivities {
	return &UserActivities{}
}

// SendWelcomeEmailActivity отправляет приветственное письмо пользователю
func (a *UserActivities) SendWelcomeEmailActivity(ctx context.Context, userID string, email string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending welcome email", "user_id", userID, "email", email)

	// Имитация отправки письма
	time.Sleep(2 * time.Second)

	logger.Info("Welcome email sent successfully", "user_id", userID)
	return nil
}

// CreateUserProfileActivity создает профиль пользователя
func (a *UserActivities) CreateUserProfileActivity(ctx context.Context, userID string, userData map[string]interface{}) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Creating user profile", "user_id", userID)

	// Имитация создания профиля
	time.Sleep(1 * time.Second)

	logger.Info("User profile created successfully", "user_id", userID)
	return nil
}

// SendNotificationActivity отправляет уведомление пользователю
func (a *UserActivities) SendNotificationActivity(ctx context.Context, userID string, message string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending notification", "user_id", userID, "message", message)

	// Имитация отправки уведомления
	time.Sleep(500 * time.Millisecond)

	logger.Info("Notification sent successfully", "user_id", userID)
	return nil
}

// ValidateUserDataActivity валидирует данные пользователя
func (a *UserActivities) ValidateUserDataActivity(ctx context.Context, userData map[string]interface{}) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Validating user data")

	// Простая валидация
	if email, ok := userData["email"].(string); !ok || email == "" {
		return fmt.Errorf("invalid email address")
	}

	if name, ok := userData["name"].(string); !ok || name == "" {
		return fmt.Errorf("invalid user name")
	}

	logger.Info("User data validation successful")
	return nil
}
