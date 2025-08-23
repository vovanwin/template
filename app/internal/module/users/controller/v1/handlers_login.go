package usersv1

import (
	"context"

	"github.com/vovanwin/platform/pkg/logger"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
)

func (i Implementation) AuthLoginPost(ctx context.Context, req *api.LoginRequest, params api.AuthLoginPostParams) (*api.AuthToken, error) {
	logger.Info(ctx, "test message")
	return &api.AuthToken{
		Access:  "111111111111",
		Refresh: "",
	}, nil
}
