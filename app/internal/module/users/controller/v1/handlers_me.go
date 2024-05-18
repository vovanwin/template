package usersv1

import (
	"context"
	"fmt"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
)

func (i Implementation) AuthMeGet(ctx context.Context, params api.AuthMeGetParams) (*api.UserMe, error) {
	user, err := i.usersService.GetMe(ctx)

	if err != nil {
		return nil, fmt.Errorf("AuthMeGet: %v", err)
	}

	return &user, nil
}
