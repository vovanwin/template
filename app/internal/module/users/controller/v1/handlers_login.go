package usersv1

import (
	"context"
	"strings"

	"github.com/google/uuid"
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

	// Генерируем пару JWT токенов
	tokenPair, err := i.jwtService.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		lg.Error(ctx, "failed to generate token pair", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Code:    500,
				Message: "Failed to generate tokens",
			},
		}
	}

	lg.Info(ctx, "successful login", zap.String("user_id", user.ID.String()))

	return &api.AuthToken{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		UserID:       user.ID,
		UserEmail:    user.Email,
	}, nil
}

// AuthRefreshPost реализует POST /auth/refresh
func (i Implementation) AuthRefreshPost(ctx context.Context, req *api.RefreshRequest, params api.AuthRefreshPostParams) (*api.AuthToken, error) {
	lg := logger.Named("users.refresh")
	lg.Info(ctx, "refresh token request")

	// Обновляем токены используя refresh токен
	tokenPair, err := i.jwtService.RefreshTokens(req.RefreshToken)
	if err != nil {
		lg.Error(ctx, "failed to refresh tokens", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 401,
			Response: api.Error{
				Code:    401,
				Message: "Invalid or expired refresh token",
			},
		}
	}

	// Извлекаем информацию о пользователе из refresh токена для ответа
	claims, err := i.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		lg.Error(ctx, "failed to validate refresh token", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 401,
			Response: api.Error{
				Code:    401,
				Message: "Invalid refresh token",
			},
		}
	}

	lg.Info(ctx, "successful token refresh", zap.String("user_id", claims.UserID))

	// Парсим UUID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		lg.Error(ctx, "failed to parse user ID", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Code:    500,
				Message: "Internal server error",
			},
		}
	}

	return &api.AuthToken{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		UserID:       userID,
		UserEmail:    claims.UserEmail,
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
