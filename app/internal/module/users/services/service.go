package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/vovanwin/template/app/internal/module/users/repository"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=usersServiceMocks

var _ UsersService = (*UsersServiceImpl)(nil)

type (
	UsersService interface {
		ValidateCredentials(ctx context.Context, email, password string) (*repository.User, error)
		GetUserByEmail(ctx context.Context, email string) (*repository.User, error)
		GetUserByID(ctx context.Context, id uuid.UUID) (*repository.User, error)
	}

	UsersServiceImpl struct {
		userRepo repository.UsersRepository
	}
)

func NewUsersServiceImpl(userRepo repository.UsersRepository) UsersService {
	return &UsersServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UsersServiceImpl) ValidateCredentials(ctx context.Context, email, password string) (*repository.User, error) {
	return s.userRepo.ValidatePassword(ctx, email, password)
}

func (s *UsersServiceImpl) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UsersServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*repository.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
