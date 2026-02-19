package auth

import (
	"context"
	authpb "github.com/vovanwin/template/pkg/auth"
)

func (s *AuthGRPCServer) OAuthURL(_ context.Context, req *authpb.OAuthURLRequest) (*authpb.OAuthURLResponse, error) {
	// TODO: implement
	panic("not implemented")
}
