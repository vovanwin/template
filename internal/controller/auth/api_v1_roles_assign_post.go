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

func (s *AuthGRPCServer) AssignRole(ctx context.Context, req *authpb.AssignRoleRequest) (*emptypb.Empty, error) {
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

	if req.GetUserId() == "" || req.GetRoleName() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and role_name are required")
	}

	targetUserID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	if err := s.authService.AssignRole(ctx, targetUserID, req.GetRoleName()); err != nil {
		s.log.Error("assign role failed", "error", err)
		return nil, status.Errorf(codes.Internal, "assign role failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
