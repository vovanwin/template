package middleware

import (
	"context"

	"github.com/vovanwin/template/app/internal/module/users/tokenDTO"
	"github.com/vovanwin/template/app/pkg/framework"
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
