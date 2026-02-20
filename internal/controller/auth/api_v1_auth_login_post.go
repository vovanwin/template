package auth

import (
	"context"

	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *AuthGRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
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

	result, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword(), ip, userAgent)
	if err != nil {
		s.log.Error("login failed", "error", err)
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}

	// Устанавливаем куку через grpc-gateway metadata
	// Кука должна быть HttpOnly, Secure, SameSite=Lax
	cookieStr := "access_token=" + result.AccessToken + "; Path=/; HttpOnly; SameSite=Lax; Max-Age=86400"
	_ = grpc.SetHeader(ctx, metadata.Pairs("Set-Cookie", cookieStr))

	// Если это HTMX запрос, добавляем заголовок перенаправления
	_ = grpc.SetHeader(ctx, metadata.Pairs("Grpc-Metadata-HX-Redirect", "/"))

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
