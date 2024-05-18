package middleware

import (
	"context"
	"github.com/vovanwin/template/internal/module/users/tokenDTO"
	"github.com/vovanwin/template/pkg/framework"
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
