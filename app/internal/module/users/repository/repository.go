package repository

import (
	"context"
	"fmt"
	"github.com/vovanwin/template/internal/middleware"
	"github.com/vovanwin/template/internal/shared/store/gen"
	users "github.com/vovanwin/template/internal/shared/store/gen/users"
	"github.com/vovanwin/template/pkg/utils"
)

var _ UsersRepo = (*EntUsersRepo)(nil)

type (
	UsersRepo interface {
		GetMe(ctx context.Context) (*gen.Users, error)
		FindForLogin(ctx context.Context, login string) (*gen.Users, error)
	}
	EntUsersRepo struct {
		ent *gen.Client
	}
)

func (r EntUsersRepo) GetMe(ctx context.Context) (*gen.Users, error) {
	claims := middleware.GetCurrentClaims(ctx)

	userMe, err := r.ent.Users.Query().Where(users.ID(claims.UserId)).First(ctx)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("repository: get user me: %v", err)
	}

	return userMe, nil
}

func (r EntUsersRepo) FindForLogin(ctx context.Context, login string) (*gen.Users, error) {
	user, err := r.ent.Users.Query().Where(users.Login(login)).First(ctx)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("repository: get user me: %v", err)
	}

	return user, nil
}

func NewEntUsersRepo(store *gen.Client) UsersRepo {
	return &EntUsersRepo{
		ent: store,
	}
}
