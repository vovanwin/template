package tokenDTO

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId   uuid.UUID `json:"user_id"`
	RoleId   int       `json:"role_id"`
	TenantId uuid.UUID `json:"tenant_id"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId    uuid.UUID `json:"user_id"`
	IsRefresh bool      `json:"is_refresh"` //TODO: рефреш токен заглушка до получаения FSSO
}
