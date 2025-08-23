package usersv1

import (
	"context"
	"time"

	"github.com/google/uuid"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
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
