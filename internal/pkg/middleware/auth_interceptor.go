package middleware

import (
	"context"
	"strings"

	"github.com/vovanwin/template/internal/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// publicMethods — методы, не требующие авторизации.
var publicMethods = map[string]bool{
	"/auth.v1.AuthService/Register":          true,
	"/auth.v1.AuthService/Login":             true,
	"/auth.v1.AuthService/RefreshToken":      true,
	"/auth.v1.AuthService/Logout":            true,
	"/auth.v1.AuthService/OAuthURL":          true,
	"/auth.v1.AuthService/OAuthCallback":     true,
	"/template.v1.TemplateService/GetHealth": true,
}

// AuthInterceptor создаёт unary gRPC interceptor для JWT авторизации.
func AuthInterceptor(jwtService jwt.JWTService, authBypass bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		if authBypass {
			// В режиме обхода авторизации прокидываем тестового пользователя
			ctx = context.WithValue(ctx, "user_id", "00000000-0000-0000-0000-000000000000")
			ctx = context.WithValue(ctx, "user_email", "dev@local")
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		if claims.TokenType != "access" {
			return nil, status.Error(codes.Unauthenticated, "not an access token")
		}

		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.UserEmail)

		return handler(ctx, req)
	}
}
