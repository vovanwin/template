package tokenDTO

import (
	"github.com/golang-jwt/jwt"
	"github.com/vovanwin/template/internal/shared/types"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId   types.UserID   `json:"user_id"`
	RoleId   int            `json:"role_id"`
	TenantId types.TenantID `json:"tenant_id"`
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
	UserId    types.UserID `json:"user_id"`
	IsRefresh bool         `json:"is_refresh"` //TODO: рефреш токен заглушка до получаения FSSO
}
