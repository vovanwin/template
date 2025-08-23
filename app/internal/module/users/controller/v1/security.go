package usersv1

import (
	"context"

	service "app/internal/module/users/services"

	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
)

var _ api.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	UsersService service.UsersService
}

func (s SecurityHandler) HandleBearerAuth(ctx context.Context, operationName string, t api.BearerAuth) (context.Context, error) {
	if operationName == "AuthLoginPost" {
		return ctx, nil
	}

	//claims, err := s.UsersService.ParseToken(t.Token)
	//if err != nil {
	//	return ctx, err
	//}

	//ctx = context.WithValue(ctx, framework.Claims, claims)

	return ctx, nil
}
