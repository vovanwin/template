package usersv1

import (
	"net/http"

	"app/config"
	service "app/internal/module/users/services"

	"github.com/go-chi/chi/v5"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
)

// Compile-time check for Handler.
var _ api.Handler = (*Implementation)(nil)

type Implementation struct {
	usersService service.UsersService
	config       *config.Config
}

func Controller(r *chi.Mux, usersService service.UsersService, config *config.Config) {
	controller := &Implementation{
		usersService: usersService,
		config:       config,
	}
	security := &SecurityHandler{
		UsersService: usersService,
	}

	srv, err := api.NewServer(
		controller,
		security,
		//api.WithTracerProvider(m.TracerProvider()),  //ждет своего часа
		//api.WithMeterProvider(m.MeterProvider()),
	)
	if err != nil {
		panic(err)
	}

	r.Mount("/", http.StripPrefix("", srv))
}
