package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims представляет claims для JWT токена
type JWTClaims struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	TokenType string `json:"token_type"` // "access" или "refresh"
	jwt.RegisteredClaims
}

// TokenPair представляет пару токенов
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// JWTService интерфейс для работы с JWT токенами
type JWTService interface {
	GenerateTokenPair(userID, userEmail string) (*TokenPair, error)
	GenerateToken(userID, userEmail string) (string, error) // оставляем для обратной совместимости
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshTokens(refreshToken string) (*TokenPair, error)
}

// DefaultJWTService реализация JWT сервиса
type DefaultJWTService struct {
	secretKey  []byte
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

// NewJWTService создает новый JWT сервис
func NewJWTService(secretKey string, tokenTTL, refreshTTL time.Duration) JWTService {
	return &DefaultJWTService{
		secretKey:  []byte(secretKey),
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
	}
}

// GenerateTokenPair создает пару access и refresh токенов
func (j *DefaultJWTService) GenerateTokenPair(userID, userEmail string) (*TokenPair, error) {
	// Генерируем access токен
	accessToken, err := j.generateTokenWithType(userID, userEmail, "access", j.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерируем refresh токен
	refreshToken, err := j.generateTokenWithType(userID, userEmail, "refresh", j.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateToken создает новый JWT токен (для обратной совместимости)
func (j *DefaultJWTService) GenerateToken(userID, userEmail string) (string, error) {
	return j.generateTokenWithType(userID, userEmail, "access", j.tokenTTL)
}

// generateTokenWithType создает токен определённого типа
func (j *DefaultJWTService) generateTokenWithType(userID, userEmail, tokenType string, ttl time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		UserEmail: userEmail,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет и парсит JWT токен
func (j *DefaultJWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// RefreshTokens обновляет токены используя refresh токен
func (j *DefaultJWTService) RefreshTokens(refreshTokenString string) (*TokenPair, error) {
	// Валидируем refresh токен
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем, что это действительно refresh токен
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Генерируем новую пару токенов
	return j.GenerateTokenPair(claims.UserID, claims.UserEmail)
}

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}

// GetUserEmailFromContext извлекает email пользователя из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	userEmail, ok := ctx.Value("user_email").(string)
	return userEmail, ok
}
