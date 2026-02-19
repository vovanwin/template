package auth

import (
	"context"

	"github.com/google/uuid"
	authpb "github.com/vovanwin/template/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthGRPCServer) CheckPermission(ctx context.Context, req *authpb.CheckPermissionRequest) (*authpb.CheckPermissionResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	allowed, err := s.authService.CheckPermission(ctx, userID, req.GetService(), req.GetAction())
	if err != nil {
		s.log.Error("check permission failed", "error", err)
		return nil, status.Errorf(codes.Internal, "check permission failed: %v", err)
	}

	return &authpb.CheckPermissionResponse{Allowed: allowed}, nil
}
