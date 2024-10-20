package tokenDTO

import (
	"app/internal/shared/types"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId   types.UserID   `json:"user_id"`
	RoleId   int            `json:"role_id"`
	TenantId types.TenantID `json:"tenant_id"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId    types.UserID `json:"user_id"`
	IsRefresh bool         `json:"is_refresh"` //TODO: рефреш токен заглушка до получаения FSSO
}
