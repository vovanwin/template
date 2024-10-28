package usersv1

import (
	api "app/internal/module/users/controller/gen"
	service "app/internal/module/users/services"
	"context"
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
