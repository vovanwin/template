package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

const (
	csrfTokenLength = 32
	csrfTokenKey    = "csrf_token"
	csrfFormField   = "csrf_token"
	csrfHeaderName  = "X-CSRF-Token"
)

// CSRFMiddleware добавляет защиту от CSRF атак
func CSRFMiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем CSRF проверку для GET, HEAD, OPTIONS запросов
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				// Генерируем или получаем CSRF токен для GET запросов
				token := getOrGenerateCSRFToken(sessionManager, r)
				w.Header().Set("X-CSRF-Token", token)
				next.ServeHTTP(w, r)
				return
			}

			// Для модифицирующих запросов проверяем CSRF токен
			if !isWebRequest(r) {
				// Для API запросов CSRF токен не обязателен (используем другие методы защиты)
				next.ServeHTTP(w, r)
				return
			}

			// Получаем токен из запроса
			var requestToken string

			// Сначала пробуем заголовок (для HTMX)
			requestToken = r.Header.Get(csrfHeaderName)

			// Если нет в заголовке, ищем в форме
			if requestToken == "" {
				requestToken = r.FormValue(csrfFormField)
			}

			// Получаем сохраненный токен из сессии
			sessionToken := sessionManager.GetString(r.Context(), csrfTokenKey)

			// Проверяем токены
			if !validateCSRFToken(requestToken, sessionToken) {
				http.Error(w, "CSRF token validation failed", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetCSRFToken возвращает CSRF токен для текущей сессии
func GetCSRFToken(sessionManager *scs.SessionManager, r *http.Request) string {
	return getOrGenerateCSRFToken(sessionManager, r)
}

// getOrGenerateCSRFToken получает существующий или генерирует новый CSRF токен
func getOrGenerateCSRFToken(sessionManager *scs.SessionManager, r *http.Request) string {
	token := sessionManager.GetString(r.Context(), csrfTokenKey)

	if token == "" {
		token = generateCSRFToken()
		sessionManager.Put(r.Context(), csrfTokenKey, token)
	}

	return token
}

// generateCSRFToken генерирует случайный CSRF токен
func generateCSRFToken() string {
	bytes := make([]byte, csrfTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		// В случае ошибки используем fallback
		return fmt.Sprintf("fallback_%d", mathrand.Int63())
	}

	return base64.URLEncoding.EncodeToString(bytes)
}

// validateCSRFToken проверяет CSRF токен с защитой от timing атак
func validateCSRFToken(requestToken, sessionToken string) bool {
	if requestToken == "" || sessionToken == "" {
		return false
	}

	// Используем constant-time сравнение для защиты от timing атак
	return subtle.ConstantTimeCompare([]byte(requestToken), []byte(sessionToken)) == 1
}
