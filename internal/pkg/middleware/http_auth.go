package middleware

import (
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

// CookieAuthMiddleware переносит JWT из куки в gRPC metadata для совместимости с AuthInterceptor.
func CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем для отладки 403
		log.Printf("[DEBUG] Request: %s %s, CSRF: %s, Auth: %s", r.Method, r.URL.Path, r.Header.Get("X-CSRF-Token"), r.Header.Get("Authorization"))

		// 1. Проверяем наличие Authorization заголовка (приоритет)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			next.ServeHTTP(w, r)
			return
		}

		// 2. Проверяем куку access_token
		cookie, err := r.Cookie("access_token")
		if err == nil && cookie.Value != "" {
			token := cookie.Value
			if !strings.HasPrefix(token, "Bearer ") {
				token = "Bearer " + token
			}

			// Устанавливаем заголовок Authorization, чтобы grpc-gateway
			// прокинул его в gRPC metadata как "authorization"
			r.Header.Set("Authorization", token)

			// Также явно добавляем в контекст для тех, кто читает из контекста напрямую
			md := metadata.Pairs("authorization", token)
			ctx := metadata.NewIncomingContext(r.Context(), md)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
