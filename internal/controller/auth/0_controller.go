package auth

import (
	"log/slog"

	"github.com/vovanwin/template/internal/service"
	authpb "github.com/vovanwin/template/pkg/auth"
	"go.uber.org/fx"
)

// Deps содержит зависимости для AuthGRPCServer.
type Deps struct {
	fx.In

	Log         *slog.Logger
	AuthService *service.AuthService
}

// AuthGRPCServer реализует gRPC сервис AuthService.
type AuthGRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	log         *slog.Logger
	authService *service.AuthService
}

// NewAuthGRPCServer создаёт новый AuthGRPCServer.
func NewAuthGRPCServer(deps Deps) *AuthGRPCServer {
	return &AuthGRPCServer{
		log:         deps.Log,
		authService: deps.AuthService,
	}
}
