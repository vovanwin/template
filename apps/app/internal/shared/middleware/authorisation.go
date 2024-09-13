package middleware

import (
	"app/pkg/framework"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInvalidAuthHeader = fmt.Errorf("invalid auth header")
)

// Мидлваре HTTP, устанавливающее значение в контексте запроса
func (m Middleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(framework.Authorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
