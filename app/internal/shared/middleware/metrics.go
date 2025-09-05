package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"

	"github.com/vovanwin/template/app/pkg/metrics"
)

// Метрика по просмотру количества запросов.
func MetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				metrics.ObserveRequest(time.Since(start), http.StatusOK)
			}()

			next.ServeHTTP(w, r)
		},
	)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	meter := otel.Meter("http-metrics")
	requestCount, _ := meter.Int64Counter("http_server_requests", metric.WithDescription("Count of HTTP requests"))
	requestDuration, _ := meter.Float64Histogram("http_server_duration", metric.WithDescription("Duration of HTTP requests"))

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Записываем количество запросов
			requestCount.Add(
				r.Context(), 1, metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", r.URL.Path),
					attribute.String("status_code", "0"), // Обновим после ответа
				),
			)

			// Засекаем время обработки запроса
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start).Seconds()

			// Обновляем метрику продолжительности
			requestDuration.Record(
				r.Context(), duration, metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", r.URL.Path),
					attribute.String("status_code", w.Header().Get("Status")),
				),
			)
		},
	)
}

// Middleware для добавления OpenTelemetry трейсов
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r.WithContext(r.Context()))
				return
			}

			//claims := tokenDTO.GetCurrentClaims(r.Context())

			tracer := otel.Tracer("")

			// Начало нового span
			ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path)
			defer span.End()

			// Пример получения информации о пользователе из контекста
			// (это зависит от вашей системы аутентификации и авторизации)
			//if claims != nil {
			//	span.SetAttributes(
			//		attribute.String("user.id", claims.UserId.String()),
			//		attribute.String("user.name", claims.TenantId.String()),
			//	)
			//}

			// Добавление информации о запросе
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
			)

			// Пример обработки ошибок
			defer func() {
				if rec := recover(); rec != nil {
					span.SetStatus(codes.Error, "panic occurred")
					span.RecordError(rec.(error))
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			// Передача контекста с span дальше по цепочке middleware
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
