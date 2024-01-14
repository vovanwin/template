package user

import "github.com/golang-jwt/jwt"

// Слой ДТО для сервиса
type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type Refresh struct {
	Refresh string `json:"refresh"`
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"user_id"`
	RoleId int    `json:"role_id"`
	Tenant string `json:"tenant"`
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
	UserId    int  `json:"user_id"`
	IsRefresh bool `json:"is_refresh"`
}
