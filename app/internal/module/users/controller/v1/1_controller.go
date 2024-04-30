package usersv1

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
	"net/http"
	"time"
)

// Compile-time check for Handler.
var _ api.Handler = (*Implementation)(nil)

type Implementation struct {
}

func Controller(r *chi.Mux) {
	controller := &Implementation{}
	security := &SecurityHandler{}

	srv, err := api.NewServer(
		controller,
		security,
		//api.WithTracerProvider(m.TracerProvider()),  //ждет своего часа
		//api.WithMeterProvider(m.MeterProvider()),
	)
	if err != nil {
		panic(err)
	}

	r.Mount("/api/v1/", http.StripPrefix("/api/v1", srv))
}

func (i Implementation) AuthMeGet(ctx context.Context, params api.AuthMeGetParams) (*api.UserMe, error) {
	return &api.UserMe{
		ID:         uuid.UUID{},
		Email:      "",
		Role:       "",
		Tenant:     "",
		CreatedAt:  time.Time{},
		Settings:   "",
		Components: nil,
	}, nil
}
