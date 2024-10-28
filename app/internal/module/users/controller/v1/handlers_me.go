package usersv1

import (
	api "app/internal/module/users/controller/gen"
	"context"
	"github.com/google/uuid"
	"time"
)

func (i Implementation) AuthMeGet(ctx context.Context, params api.AuthMeGetParams) (*api.UserMe, error) {

	return &api.UserMe{
		ID:    uuid.UUID{},
		Email: "ssssssssss",
		Role: api.OptString{
			Value: "",
			Set:   true,
		},
		Tenant:     "",
		CreatedAt:  time.Time{},
		Settings:   "",
		Components: nil,
	}, nil
}
