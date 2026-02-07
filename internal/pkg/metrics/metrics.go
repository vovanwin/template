package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler возвращает HTTP handler для Prometheus /metrics endpoint.
// Метрики собираются автоматически через OTEL Prometheus exporter,
// который регистрирует prometheus.DefaultRegisterer при инициализации.
func Handler() http.Handler {
	return promhttp.Handler()
}
