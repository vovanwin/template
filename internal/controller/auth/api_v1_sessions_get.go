package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/vovanwin/template/internal/pkg/jwt"
	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthGRPCServer) ListSessions(ctx context.Context, _ *emptypb.Empty) (*authpb.SessionsResponse, error) {
	userIDStr, ok := jwt.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}

	sessions, err := s.authService.ListSessions(ctx, userID)
	if err != nil {
		s.log.Error("list sessions failed", "error", err)
		return nil, status.Errorf(codes.Internal, "list sessions failed: %v", err)
	}

	pbSessions := make([]*authpb.SessionInfo, 0, len(sessions))
	for _, sess := range sessions {
		pbSessions = append(pbSessions, &authpb.SessionInfo{
			Id:        sess.ID.String(),
			Ip:        sess.IP,
			UserAgent: sess.UserAgent,
			CreatedAt: timestamppb.New(sess.CreatedAt),
			ExpiresAt: timestamppb.New(sess.ExpiresAt),
		})
	}

	return &authpb.SessionsResponse{Sessions: pbSessions}, nil
}
