package usersv1

import (
	"context"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
	service "github.com/vovanwin/template/internal/module/users/services"
	"github.com/vovanwin/template/pkg/framework"
)

var _ api.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	UsersService service.UsersService
}

func (s SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	claims, err := s.UsersService.ParseToken(t.Token)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, framework.Claims, claims)

	return ctx, nil
}
