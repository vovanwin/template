package middleware

import (
	"app/internal/module/users/tokenDTO"
	"app/pkg/framework"
	"context"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func GetCurrentClaims(ctx context.Context) *tokenDTO.TokenClaims {
	claims := ctx.Value(framework.Claims).(*tokenDTO.TokenClaims)
	return claims
}
