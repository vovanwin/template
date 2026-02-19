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

func (s *AuthGRPCServer) GetProfile(ctx context.Context, _ *emptypb.Empty) (*authpb.ProfileResponse, error) {
	userIDStr, ok := jwt.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}

	profile, err := s.authService.GetProfile(ctx, userID)
	if err != nil {
		s.log.Error("get profile failed", "error", err)
		return nil, status.Errorf(codes.Internal, "get profile failed: %v", err)
	}

	return &authpb.ProfileResponse{
		Id:        profile.ID.String(),
		Email:     profile.Email,
		Name:      profile.Name,
		AvatarUrl: profile.AvatarURL,
		IsActive:  profile.IsActive,
		Roles:     []string{profile.Role},
		CreatedAt: timestamppb.New(profile.CreatedAt),
		UpdatedAt: timestamppb.New(profile.UpdatedAt),
	}, nil
}
