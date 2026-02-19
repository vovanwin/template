package auth

import (
	"context"
	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthGRPCServer) GetLinkedAccounts(_ context.Context, req *emptypb.Empty) (*authpb.LinkedAccountsResponse, error) {
	// TODO: implement
	panic("not implemented")
}
