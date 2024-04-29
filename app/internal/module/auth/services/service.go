package service

import (
	"context"
	"github.com/vovanwin/template/internal/module/auth/repository"
)

var _ AuthService = (*AuthServiceImpl)(nil)

type (
	AuthService interface {
		GetLogin(ctx context.Context) error
	}
	AuthServiceImpl struct {
		AuthRepo repository.AuthRepo
	}
)

func NewAuthServiceImpl(AuthRepo repository.AuthRepo) AuthService {
	return &AuthServiceImpl{
		AuthRepo: AuthRepo,
	}
}
func (r AuthServiceImpl) GetLogin(ctx context.Context) error {
	err := r.AuthRepo.GetLogin(ctx)
	return err
}
