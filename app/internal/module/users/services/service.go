package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/vovanwin/template/config"
	api "github.com/vovanwin/template/internal/module/users/controller/gen"
	"github.com/vovanwin/template/internal/module/users/repository"
	"github.com/vovanwin/template/internal/module/users/tokenDTO"
	"github.com/vovanwin/template/internal/shared/store/gen"
	"github.com/vovanwin/template/pkg/utils"
	"github.com/vovanwin/template/pkg/utils/hasher"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=usersServiceMocks

var _ UsersService = (*UsersServiceImpl)(nil)

type (
	UsersService interface {
		GetMe(ctx context.Context) (api.UserMe, error)
		GetTokens(ctx context.Context, req *api.LoginRequest) (*api.AuthToken, error)
		ParseToken(accessToken string) (*tokenDTO.TokenClaims, error)
	}
	UsersServiceImpl struct {
		UsersRepo repository.UsersRepo
		config    *config.Config
	}
)

func NewUsersServiceImpl(UsersRepo repository.UsersRepo, config *config.Config) UsersService {
	return &UsersServiceImpl{
		UsersRepo: UsersRepo,
		config:    config,
	}
}

func (u UsersServiceImpl) GetMe(ctx context.Context) (api.UserMe, error) {
	user, err := u.UsersRepo.GetMe(ctx)
	if err != nil {
		return api.UserMe{}, err
	}
	userResponse := api.UserMe{
		ID:         uuid.UUID(user.ID),
		Email:      user.Login,
		Role:       "Админ",
		Tenant:     "asd",
		CreatedAt:  user.CreatedAt,
		Settings:   "",
		Components: nil,
	}
	return userResponse, nil
}

func (u UsersServiceImpl) GetTokens(ctx context.Context, req *api.LoginRequest) (*api.AuthToken, error) {
	user, err := u.UsersRepo.FindForLogin(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("GetTokens GetMe: %v", err)
	}

	isValid, err := hasher.ComparePasswordAndHash(req.Password, user.Password)
	if isValid == false {
		return nil, utils.ErrValidPassword
	}

	// сгенерировать токен
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("GetTokens generateAccessToken: %v", err)
	}
	// сгенерировать токен
	refreshToken, err := u.generateRefreshToken(user)

	return &api.AuthToken{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (u UsersServiceImpl) generateAccessToken(user *gen.Users) (string, error) {
	// сгенерировать токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenDTO.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(u.config.JWT.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.ID,
	})

	// подписать токен
	tokenString, err := token.SignedString([]byte(u.config.SignKey))
	if err != nil {
		return "", utils.ErrCannotSignToken
	}

	return tokenString, nil
}

// GenerateRefreshToken TODO:заглушка для refresh token до получения FSSO, потом токен будет запрашивать у него и сами будет генерировать только access токен
func (u *UsersServiceImpl) generateRefreshToken(user *gen.Users) (string, error) {
	// сгенерировать токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenDTO.RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(u.config.JWT.RefreshTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:    user.ID,
		IsRefresh: true,
	})

	// подписать токен
	tokenString, err := token.SignedString([]byte(u.config.SignKey))
	if err != nil {
		return "", utils.ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *UsersServiceImpl) ParseToken(accessToken string) (*tokenDTO.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenDTO.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.config.SignKey), nil
	})

	if err != nil {
		return nil, utils.ErrCannotParseToken
	}

	claims, ok := token.Claims.(*tokenDTO.TokenClaims)
	if !ok {
		return nil, utils.ErrCannotParseToken
	}

	return claims, nil
}
