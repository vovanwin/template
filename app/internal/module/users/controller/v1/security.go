package usersv1

import (
	"context"
	"fmt"
	"strings"

	service "github.com/vovanwin/template/app/internal/module/users/services"
	"github.com/vovanwin/template/app/pkg/jwt"

	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
)

var _ api.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	UsersService service.UsersService
	JWTService   jwt.JWTService
}

func (s SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	fmt.Printf("DEBUG SecurityHandler: operationName='%s', token='%s'\n", operationName, t.Token)

	// Login не требует токена - он должен обрабатываться без security
	if operationName == "AuthLoginPost" {
		fmt.Printf("DEBUG SecurityHandler: AuthLoginPost detected, this should not happen!\n")
		return ctx, nil
	}

	// Извлекаем токен, удаляя "Bearer " префикс если есть
	tokenString := t.Token
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	fmt.Printf("DEBUG SecurityHandler: validating token='%s'\n", tokenString)

	// Валидируем JWT токен
	claims, err := s.JWTService.ValidateToken(tokenString)
	if err != nil {
		fmt.Printf("DEBUG SecurityHandler: token validation failed: %v\n", err)
		return ctx, fmt.Errorf("invalid token: %w", err)
	}

	fmt.Printf("DEBUG SecurityHandler: token valid, userID='%s', userEmail='%s'\n", claims.UserID, claims.UserEmail)

	// Добавляем данные пользователя в контекст
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_email", claims.UserEmail)

	return ctx, nil
}
