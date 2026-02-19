package auth

import (
	"context"
	authpb "github.com/vovanwin/template/pkg/auth"
)

func (s *AuthGRPCServer) OAuthCallback(_ context.Context, req *authpb.OAuthCallbackRequest) (*authpb.AuthResponse, error) {
	// TODO: implement
	panic("not implemented")
}
