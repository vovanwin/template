package auth

import (
	"context"

	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthGRPCServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	result, err := s.authService.Register(ctx, req.GetEmail(), req.GetPassword(), req.GetName())
	if err != nil {
		s.log.Error("register failed", "error", err)
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}

	return &authpb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &authpb.UserInfo{
			Id:    result.User.ID.String(),
			Email: result.User.Email,
			Name:  result.User.FirstName,
		},
	}, nil
}
