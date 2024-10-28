package usersv1

import (
	"app/config"
	api "app/internal/module/users/controller/gen"
	service "app/internal/module/users/services"
	"github.com/go-chi/chi/v5"
	"net/http"
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

	r.Mount("/api/v1/", http.StripPrefix("/api/v1", srv))
}
