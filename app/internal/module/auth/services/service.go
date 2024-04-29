package service

import (
	"context"
	"github.com/vovanwin/template/internal/module/auth/repository"
	"github.com/vovanwin/template/internal/shared/types"
	"time"
)

var _ AuthService = (*AuthServiceImpl)(nil)

type (
	AuthService interface {
		GetLogin(ctx context.Context) ([]Data, error)
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

func (r AuthServiceImpl) GetLogin(ctx context.Context) ([]Data, error) {
	users, err := r.AuthRepo.GetLogin(ctx)
	usersResponse := make([]Data, 0, len(users))
	for _, m := range users {
		usersResponse = append(usersResponse, Data{
			ID:        m.ID,
			Email:     m.Email,
			DeletedAt: m.DeletedAt,
			UpdatedAt: m.UpdatedAt,
			CreatedAt: m.CreatedAt,
		})
	}
	return usersResponse, err
}

type Data struct {
	ID        types.UserID `json:"id,omitempty"`
	Email     string       `json:"email,omitempty"`
	DeletedAt time.Time    `json:"deleted_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
}
