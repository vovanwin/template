package auth

import (
	"context"

	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *AuthGRPCServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.AuthResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	ip := ""
	if p, ok := peer.FromContext(ctx); ok {
		ip = p.Addr.String()
	}

	userAgent := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			ip = xff[0]
		}
	}

	result, err := s.authService.RefreshToken(ctx, req.GetRefreshToken(), ip, userAgent)
	if err != nil {
		s.log.Error("refresh token failed", "error", err)
		return nil, status.Errorf(codes.Unauthenticated, "refresh failed: %v", err)
	}

	return &authpb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User: &authpb.UserInfo{
			Id:    result.User.ID.String(),
			Email: result.User.Email,
			Name:  result.User.Name,
		},
	}, nil
}
