package repository

import (
	"context"
	"github.com/vovanwin/template/internal/shared/store/gen"
)

var _ AuthRepo = (*EntAuthRepo)(nil)

type (
	AuthRepo interface {
		GetLogin(ctx context.Context) ([]*gen.User, error)
	}
	EntAuthRepo struct {
		ent *gen.Client
	}
)

func (r EntAuthRepo) GetLogin(ctx context.Context) ([]*gen.User, error) {
	users := r.ent.User.Query().AllX(ctx)

	return users, nil
}

func NewEntAuthRepo(store *gen.Client) AuthRepo {
	return &EntAuthRepo{
		ent: store,
	}
}
