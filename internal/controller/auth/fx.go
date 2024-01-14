package auth

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	authSer "template/internal/domain/user/service/auth"
)

type UserController struct {
	AuthService authSer.AuthService
	logger      *slog.Logger
}

func InitIndexController(r *chi.Mux, AuthService authSer.AuthService, logger *slog.Logger) {
	controller := &UserController{
		AuthService: AuthService,
		logger:      logger,
	}

	r.Route("/api/v1/auth", func(r chi.Router) {
		// Public Routes
		r.Post("/login", controller.login)
		r.Post("/refresh", controller.login)

	})
}
