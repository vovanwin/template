package user

import (
	"context"
	"template/internal/entity"
	"template/internal/repository/user"
	"time"
)

type (
	UserService interface {
		GetByID(ctx context.Context, id int64) (user entity.User, err error)
		GetByLogin(ctx context.Context, login string) (user entity.User, err error)
		Delete(ctx context.Context, id int64) (err error)
	}
	UserImpl struct {
		repo           user.UserRepo
		contextTimeout time.Duration
	}
)

func NewUserImpl(repo user.UserRepo, timeout time.Duration) UserService {
	if repo == nil {
		panic("user Repository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &UserImpl{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (c *UserImpl) GetByID(ctx context.Context, id int64) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	user, err = c.repo.GetByID(ctx, id)
	return user, err
}

func (c *UserImpl) GetByLogin(ctx context.Context, login string) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	user, err = c.repo.GetByLogin(ctx, login)
	return user, err
}

func (c *UserImpl) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.repo.Delete(ctx, id)
	return err
}
