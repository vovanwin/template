package auth

import (
	"context"
	authpb "github.com/vovanwin/template/pkg/auth"
)

func (s *AuthGRPCServer) LinkOAuth(_ context.Context, req *authpb.OAuthCallbackRequest) (*authpb.LinkResponse, error) {
	// TODO: implement
	panic("not implemented")
}
