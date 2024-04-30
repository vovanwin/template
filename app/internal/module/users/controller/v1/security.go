package usersv1

import (
	"context"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
)

var _ api.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
}

func (s SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	// тут проверка токена (валидоность и что нибудь ещё по необходимости)
	return ctx, nil
}
