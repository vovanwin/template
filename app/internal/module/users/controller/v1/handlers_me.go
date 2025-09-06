package usersv1

import (
	"context"

	"github.com/vovanwin/platform/pkg/logger"
	"github.com/vovanwin/template/app/pkg/jwt"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
	"go.uber.org/zap"
)

func (i Implementation) AuthMeGet(ctx context.Context, params api.AuthMeGetParams) (*api.UserMe, error) {
	lg := logger.Named("users.me")
	lg.Info(ctx, "me request")

	// Получаем данные пользователя из контекста (установленные SecurityHandler)
	userEmail, ok := jwt.GetUserEmailFromContext(ctx)
	if !ok {
		lg.Error(ctx, "user not authenticated")
		return nil, &api.ErrorStatusCode{
			StatusCode: 401,
			Response: api.Error{
				Code:    401,
				Message: "User not authenticated",
			},
		}
	}

	// Получаем полную информацию о пользователе из базы
	user, err := i.usersService.GetUserByEmail(ctx, userEmail)
	if err != nil || user == nil {
		lg.Error(ctx, "user not found", zap.String("email", userEmail), zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 404,
			Response: api.Error{
				Code:    404,
				Message: "User not found",
			},
		}
	}

	// Формируем ответ
	userMe := &api.UserMe{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Settings:  user.Settings,
	}

	// Опциональные поля
	if user.Role != "" {
		userMe.Role = api.OptString{Value: user.Role, Set: true}
	}

	if user.TenantID != nil {
		userMe.Tenant = *user.TenantID
	}

	// Компоненты (пока захардкодим как в комментарии к схеме)
	userMe.Components = []string{"/monitoringmap"}

	lg.Info(ctx, "successful me request", zap.String("user_id", user.ID.String()))

	return userMe, nil
}
