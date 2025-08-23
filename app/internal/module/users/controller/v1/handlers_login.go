package usersv1

import (
	"context"

	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
)

func (i Implementation) AuthLoginPost(ctx context.Context, req *api.LoginRequest, params api.AuthLoginPostParams) (*api.AuthToken, error) {

	return &api.AuthToken{
		Access:  "111111111111",
		Refresh: "",
	}, nil
}
