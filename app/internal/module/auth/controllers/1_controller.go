package controller

import (
	"github.com/go-chi/chi/v5"
	service "github.com/vovanwin/template/internal/module/auth/services"
)

type Implementation struct {
	service service.AuthService
}

func Controller(
	r *chi.Mux,
	service service.AuthService,
) {
	controller := &Implementation{
		service: service,
	}

	r.Group(func(r chi.Router) {

		r.Route("/api/v1/test", func(r chi.Router) {
			r.Get("/", controller.TEST)

		})
	})
}
