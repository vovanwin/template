package user

import (
	"context"
	"errors"
	"template/internal/domain/user/entity"
	userRep "template/internal/domain/user/repository/user"
	"time"
)

type (
	UserService interface {
		GetByID(ctx context.Context, id int) (user entity.User, err error)
		GetByLogin(ctx context.Context, login string) (user entity.User, err error)
		Delete(ctx context.Context, id int) (err error)
	}
	UserImpl struct {
		userRep        userRep.UserRepo
		contextTimeout time.Duration
	}
)

func NewUserImpl(userRep userRep.UserRepo, timeout time.Duration) UserService {
	if userRep == nil {
		panic("user Repository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &UserImpl{
		userRep:        userRep,
		contextTimeout: timeout,
	}
}

func (c *UserImpl) GetByID(ctx context.Context, id int) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	user, err = c.userRep.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, userRep.ErrNotFound) {
			return user, ErrUserNotFound
		}
		return user, err
	}
	return user, err
}

func (c *UserImpl) GetByLogin(ctx context.Context, login string) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	user, err = c.userRep.GetByLogin(ctx, login)
	return user, err
}

func (c *UserImpl) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	err := c.userRep.Delete(ctx, id)
	return err
}
