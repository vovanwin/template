package otel

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPMiddleware возвращает chi-совместимый middleware, оборачивающий каждый запрос
// в OTEL span и записывающий метрики (duration, status code).
func HTTPMiddleware(operation string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewMiddleware(operation)(next)
	}
}
