package usersv1

import (
	"context"
	"fmt"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
)

func (i Implementation) AuthLoginPost(ctx context.Context, req *api.LoginRequest, params api.AuthLoginPostParams) (*api.AuthToken, error) {
	token, err := i.usersService.GetTokens(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("AuthLoginPost: %v", err)
	}

	return &api.AuthToken{
		Access:  token.Access,
		Refresh: token.Refresh,
	}, nil
}
