package auth

import (
	"context"
	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthGRPCServer) UnlinkOAuth(_ context.Context, req *authpb.UnlinkRequest) (*emptypb.Empty, error) {
	// TODO: implement
	panic("not implemented")
}
