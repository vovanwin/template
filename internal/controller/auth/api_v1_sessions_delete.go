package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/vovanwin/template/internal/pkg/jwt"
	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthGRPCServer) RevokeSession(ctx context.Context, req *authpb.RevokeSessionRequest) (*emptypb.Empty, error) {
	userIDStr, ok := jwt.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}

	sessionID, err := uuid.Parse(req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid session_id")
	}

	if err := s.authService.RevokeSession(ctx, userID, sessionID); err != nil {
		s.log.Error("revoke session failed", "error", err)
		return nil, status.Errorf(codes.Internal, "revoke session failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
