package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/vovanwin/template/internal/pkg/events"
	"github.com/vovanwin/template/internal/pkg/jwt"
	"github.com/vovanwin/template/internal/pkg/utils/hasher"
	"github.com/vovanwin/template/internal/repository"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *repository.User
}

type Profile struct {
	ID        uuid.UUID
	Email     string
	Name      string
	AvatarURL string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionInfo struct {
	ID        uuid.UUID
	IP        string
	UserAgent string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type LinkedAccount struct {
	Provider   string
	ProviderID string
	Email      string
	CreatedAt  time.Time
}

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
	jwt         jwt.JWTService
	bus         *events.Bus
	log         *slog.Logger
}

func NewAuthService(
	userRepo *repository.UserRepo,
	sessionRepo *repository.SessionRepo,
	jwtService jwt.JWTService,
	bus *events.Bus,
	log *slog.Logger,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwt:         jwtService,
		bus:         bus,
		log:         log,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*AuthResult, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	passwordHash, err := hasher.CreateHash(password, hasher.DefaultParams)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, email, passwordHash, name)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	refreshHash := hashToken(tokens.RefreshToken)
	_, err = s.sessionRepo.Create(ctx, user.ID, refreshHash, "", "", time.Now().Add(30*24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &AuthResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password, ip, userAgent string) (*AuthResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("user is deactivated")
	}

	match, err := hasher.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("compare password: %w", err)
	}
	if !match {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	refreshHash := hashToken(tokens.RefreshToken)
	_, err = s.sessionRepo.Create(ctx, user.ID, refreshHash, ip, userAgent, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	// Отправляем уведомление пользователю
	_ = s.bus.Publish(events.Event{
		UserID:  user.ID.String(),
		Message: fmt.Sprintf("👋 Привет, %s! Успешный вход в систему.", user.FirstName),
		Type:    "success",
	})

	return &AuthResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	refreshHash := hashToken(refreshToken)
	return s.sessionRepo.DeleteByTokenHash(ctx, refreshHash)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, ip, userAgent string) (*AuthResult, error) {
	claims, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	refreshHash := hashToken(refreshToken)
	session, err := s.sessionRepo.GetByTokenHash(ctx, refreshHash)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Delete old session
	if err := s.sessionRepo.DeleteByTokenHash(ctx, refreshHash); err != nil {
		return nil, fmt.Errorf("delete old session: %w", err)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new token pair
	tokens, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	newRefreshHash := hashToken(tokens.RefreshToken)
	_, err = s.sessionRepo.Create(ctx, user.ID, newRefreshHash, ip, userAgent, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("create new session: %w", err)
	}

	return &AuthResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &Profile{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.FirstName,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, avatarURL string) (*Profile, error) {
	fields := make(map[string]interface{})
	if name != "" {
		fields["name"] = name
	}
	if avatarURL != "" {
		fields["avatar_url"] = avatarURL
	}
	if len(fields) == 0 {
		return s.GetProfile(ctx, userID)
	}

	user, err := s.userRepo.Update(ctx, userID, fields)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &Profile{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.FirstName,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) ListSessions(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error) {
	sessions, err := s.sessionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	result := make([]SessionInfo, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, SessionInfo{
			ID:        sess.ID,
			IP:        sess.IP,
			UserAgent: sess.UserAgent,
			CreatedAt: sess.CreatedAt,
			ExpiresAt: sess.ExpiresAt,
		})
	}
	return result, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	// Verify session belongs to user by listing and checking
	sessions, err := s.sessionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("list sessions: %w", err)
	}

	found := false
	for _, sess := range sessions {
		if sess.ID == sessionID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("session not found")
	}

	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *AuthService) AssignRole(ctx context.Context, targetUserID uuid.UUID, roleName string) error {
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	_, err = s.userRepo.SetRole(ctx, targetUserID, roleName)
	if err != nil {
		return fmt.Errorf("set role: %w", err)
	}
	return nil
}

func (s *AuthService) RemoveRole(ctx context.Context, targetUserID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	_, err = s.userRepo.SetRole(ctx, targetUserID, "user")
	if err != nil {
		return fmt.Errorf("remove role: %w", err)
	}
	return nil
}

// CheckPermission проверяет, имеет ли пользователь доступ к действию.
// Простая RBAC: admin может всё, user — ограниченный набор действий.
func (s *AuthService) CheckPermission(ctx context.Context, targetUserID uuid.UUID, svc, action string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return false, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return false, fmt.Errorf("user not found")
	}

	if user.Role == "admin" {
		return true, nil
	}

	// user роль — разрешаем read-only действия
	switch action {
	case "read", "list", "view":
		return true, nil
	default:
		return false, nil
	}
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
