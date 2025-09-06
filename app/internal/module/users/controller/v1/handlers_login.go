package usersv1

import (
	"context"
	"strings"

	"github.com/vovanwin/platform/pkg/logger"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
	"go.uber.org/zap"
)

func (i Implementation) AuthLoginPost(ctx context.Context, req *api.LoginRequest, params api.AuthLoginPostParams) (*api.AuthToken, error) {
	lg := logger.Named("users.login")
	lg.Info(ctx, "login request", zap.String("email", req.Email))

	// Очищаем email
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Валидируем учетные данные
	user, err := i.usersService.ValidateCredentials(ctx, email, req.Password)
	if err != nil {
		lg.Error(ctx, "failed to validate credentials", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Code:    500,
				Message: "Internal server error",
			},
		}
	}

	if user == nil {
		lg.Warn(ctx, "invalid credentials", zap.String("email", email))
		return nil, &api.ErrorStatusCode{
			StatusCode: 401,
			Response: api.Error{
				Code:    401,
				Message: "Invalid email or password",
			},
		}
	}

	// Генерируем JWT токен
	token, err := i.jwtService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		lg.Error(ctx, "failed to generate token", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Code:    500,
				Message: "Failed to generate token",
			},
		}
	}

	lg.Info(ctx, "successful login", zap.String("user_id", user.ID.String()))

	return &api.AuthToken{
		Token:     token,
		UserID:    user.ID,
		UserEmail: user.Email,
	}, nil
}

// AuthLogoutPost реализует POST /auth/logout
func (i Implementation) AuthLogoutPost(ctx context.Context, params api.AuthLogoutPostParams) (*api.LogoutResponse, error) {
	lg := logger.Named("users.logout")
	lg.Info(ctx, "logout request")

	// В JWT logout просто означает, что клиент должен удалить токен
	// На сервере ничего делать не нужно (токен истечет сам)
	return &api.LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}
