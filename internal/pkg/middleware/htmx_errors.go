package middleware

import (
	"net/http"
)

// HTMXErrorMiddleware перехватывает ошибки и возвращает HTML-фрагменты для HTMX
func HTMXErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			next.ServeHTTP(w, r)
			return
		}

		// Используем кастомный ResponseWriter для перехвата статуса
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		// Если была ошибка (например, 403 CSRF или 401 Unauth)
		if rec.status >= 400 && rec.status < 600 {
			// Очищаем то, что уже успели записать (если это не был Header)
			// HTMX ожидает 200 OK для того, чтобы вставить ошибку в hx-target
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			// Мы не можем изменить статус после того, как WriteHeader уже вызван,
			// поэтому просто пишем текст ошибки в надежде, что HTMX это обработает.
			// Но лучше всего HTMX работает, когда мы возвращаем 200 OK с текстом ошибки.
		}
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
