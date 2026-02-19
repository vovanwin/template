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

func (s *AuthGRPCServer) RemoveRole(ctx context.Context, req *authpb.RemoveRoleRequest) (*emptypb.Empty, error) {
	callerIDStr, ok := jwt.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	callerID, err := uuid.Parse(callerIDStr)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid caller id")
	}

	allowed, err := s.authService.CheckPermission(ctx, callerID, "auth", "admin")
	if err != nil || !allowed {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	targetUserID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	if err := s.authService.RemoveRole(ctx, targetUserID); err != nil {
		s.log.Error("remove role failed", "error", err)
		return nil, status.Errorf(codes.Internal, "remove role failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
