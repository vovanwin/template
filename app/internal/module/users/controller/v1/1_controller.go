package usersv1

import (
	"net/http"

	"github.com/vovanwin/template/app/config"
	service "github.com/vovanwin/template/app/internal/module/users/services"
	"github.com/vovanwin/template/app/pkg/jwt"

	"github.com/go-chi/chi/v5"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
)

// Compile-time check for Handler.
var _ api.Handler = (*Implementation)(nil)

type Implementation struct {
	usersService service.UsersService
	config       *config.Config
	jwtService   jwt.JWTService
}

func Controller(r *chi.Mux, usersService service.UsersService, config *config.Config, jwtService jwt.JWTService) {
	controller := &Implementation{
		usersService: usersService,
		config:       config,
		jwtService:   jwtService,
	}
	security := &SecurityHandler{
		UsersService: usersService,
		JWTService:   jwtService,
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
