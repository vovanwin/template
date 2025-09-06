package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// SessionSecurityMiddleware добавляет дополнительную безопасность для сессий
func SessionSecurityMiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверка на session fixation атаки
			if sessionManager.Exists(r.Context(), "user_id") {
				// Проверяем, не истекла ли сессия по времени бездействия
				lastActivity := sessionManager.GetTime(r.Context(), "last_activity")
				if !lastActivity.IsZero() && time.Since(lastActivity) > 30*time.Minute {
					// Очищаем истекшую сессию
					sessionManager.Clear(r.Context())
					sessionManager.Put(r.Context(), "flash_message", "Session expired due to inactivity")

					// Если это API запрос, возвращаем 401
					if isAPIRequest(r) {
						http.Error(w, "Session expired", http.StatusUnauthorized)
						return
					}

					// Для веб-запросов перенаправляем на логин
					if isWebRequest(r) {
						http.Redirect(w, r, "/web/login", http.StatusSeeOther)
						return
					}
				}

				// Обновляем время последней активности
				sessionManager.Put(r.Context(), "last_activity", time.Now())

				// Проверка IP адреса (опционально)
				sessionIP := sessionManager.GetString(r.Context(), "session_ip")
				currentIP := getClientIP(r)

				if sessionIP != "" && sessionIP != currentIP {
					slog.Warn("Session IP mismatch",
						"session_ip", sessionIP,
						"current_ip", currentIP,
						"user_agent", r.UserAgent())

					// В production режиме можно очистить сессию
					// sessionManager.Clear(r.Context())
					// http.Redirect(w, r, "/web/login", http.StatusSeeOther)
					// return
				}
			}

			// Добавляем заголовки безопасности
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			next.ServeHTTP(w, r)
		})
	}
}

// AuthRequiredMiddleware проверяет, что пользователь авторизован
func AuthRequiredMiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sessionManager.Exists(r.Context(), "user_id") {
				if isAPIRequest(r) {
					http.Error(w, "Authentication required", http.StatusUnauthorized)
					return
				}

				if isWebRequest(r) {
					// Сохраняем URL, куда пользователь хотел попасть
					sessionManager.Put(r.Context(), "redirect_after_login", r.URL.Path)
					http.Redirect(w, r, "/web/login", http.StatusSeeOther)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Вспомогательные функции
func isAPIRequest(r *http.Request) bool {
	return r.Header.Get("Accept") == "application/json" ||
		r.Header.Get("Content-Type") == "application/json"
}

func isWebRequest(r *http.Request) bool {
	return r.URL.Path[0:4] == "/web"
}

// GetClientIP возвращает IP адрес клиента (экспортируемая функция)
func GetClientIP(r *http.Request) string {
	return getClientIP(r)
}

func getClientIP(r *http.Request) string {
	// Проверяем заголовки прокси
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	return r.RemoteAddr
}
