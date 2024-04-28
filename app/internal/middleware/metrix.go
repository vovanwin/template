package middleware

import (
	"github.com/vovanwin/template/pkg/metrics"
	"net/http"
	"time"
)

// Метрика по просмотру количества запросов
func MetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			metrics.ObserveRequest(time.Since(start), http.StatusOK)
		}()

		next.ServeHTTP(w, r)
	})
}
