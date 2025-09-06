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
	jwt.RegisteredClaims
}

// JWTService интерфейс для работы с JWT токенами
type JWTService interface {
	GenerateToken(userID, userEmail string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// DefaultJWTService реализация JWT сервиса
type DefaultJWTService struct {
	secretKey []byte
	tokenTTL  time.Duration
}

// NewJWTService создает новый JWT сервис
func NewJWTService(secretKey string, tokenTTL time.Duration) JWTService {
	return &DefaultJWTService{
		secretKey: []byte(secretKey),
		tokenTTL:  tokenTTL,
	}
}

// GenerateToken создает новый JWT токен
func (j *DefaultJWTService) GenerateToken(userID, userEmail string) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenTTL)),
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
