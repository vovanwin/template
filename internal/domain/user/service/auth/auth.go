package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"log/slog"
	"template/config"
	"template/internal/domain/user/entity"
	userRep "template/internal/domain/user/repository/user"
	"template/pkg/utils/hasher"
	"time"
)

type (
	AuthService interface {
		GetTokens(ctx context.Context, input AuthGenerateTokenInput) (Tokens, error)
		generateAccessToken(user entity.User) (string, error)
		GenerateRefreshToken(user entity.User) (string, error)
		ParseToken(accessToken string) (*TokenClaims, error)
		RefreshParseToken(refreshToken string) (*RefreshTokenClaims, error)
	}
	AuthImpl struct {
		userRep        userRep.UserRepo
		contextTimeout time.Duration
		config         config.Config
		*slog.Logger
	}
)

func NewAuthImpl(userRep userRep.UserRepo, config config.Config, timeout time.Duration, log *slog.Logger) AuthService {
	if userRep == nil {
		panic("user Repository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &AuthImpl{
		userRep:        userRep,
		contextTimeout: timeout,
		config:         config,
		Logger:         log,
	}
}
func (s *AuthImpl) GetTokens(ctx context.Context, input AuthGenerateTokenInput) (Tokens, error) {
	user, err := s.userRep.GetByLogin(ctx, input.Username)
	if err != nil {
		if errors.Is(err, userRep.ErrNotFound) {
			return Tokens{}, ErrUserNotFound
		}
		s.Logger.Error("AuthService.GenerateToken: cannot get user: %v", err)
		return Tokens{}, ErrCannotGetUser
	}

	isValid, err := hasher.ComparePasswordAndHash(input.Password, user.Password)

	if isValid == false {
		return Tokens{}, ErrValidPassword
	}

	// сгенерировать токен
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return Tokens{}, err
	}

	// сгенерировать токен
	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{Access: accessToken, Refresh: refreshToken}, nil
}

func (s *AuthImpl) generateAccessToken(user entity.User) (string, error) {
	// сгенерировать токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.JWT.AccessTtl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.ID,
		Tenant: user.Tenant,
		RoleId: user.RoleId,
	})

	// подписать токен
	tokenString, err := token.SignedString([]byte(s.config.JWT.SighKey))
	if err != nil {
		log.Error("AuthService.GenerateToken: cannot sign token: %v", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthImpl) GenerateRefreshToken(user entity.User) (string, error) {
	// сгенерировать токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.JWT.RefreshTtl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:    user.ID,
		IsRefresh: true,
	})

	// подписать токен
	tokenString, err := token.SignedString([]byte(s.config.JWT.SighKey))
	if err != nil {
		log.Error("AuthService.GenerateToken: cannot sign token: %v", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthImpl) ParseToken(accessToken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.config.JWT.SighKey), nil
	})

	if err != nil {
		return nil, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrCannotParseToken
	}

	return claims, nil
}

func (s *AuthImpl) RefreshParseToken(refreshToken string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.config.JWT.SighKey), nil
	})

	if err != nil {
		return nil, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, ErrCannotParseToken
	}

	return claims, nil
}
