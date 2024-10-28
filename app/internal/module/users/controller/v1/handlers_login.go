package usersv1

import (
	api "app/internal/module/users/controller/gen"
	"context"
)

func (i Implementation) AuthLoginPost(ctx context.Context, req *api.LoginRequest, params api.AuthLoginPostParams) (*api.AuthToken, error) {

	return &api.AuthToken{
		Access:  "111111111111",
		Refresh: "",
	}, nil
}
