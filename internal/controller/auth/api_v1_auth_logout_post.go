package auth

import (
	"context"

	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthGRPCServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*emptypb.Empty, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	if err := s.authService.Logout(ctx, req.GetRefreshToken()); err != nil {
		s.log.Error("logout failed", "error", err)
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
