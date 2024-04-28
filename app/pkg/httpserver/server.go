package httpserver

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	customMiddleware "github.com/vovanwin/template/internal/middleware"
	logMiddleware "github.com/vovanwin/template/pkg/fxslog/logger"
	"log/slog"
	"net/http"
	"time"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger            *slog.Logger  `option:"mandatory" validate:"required"`
	isProduction      bool          `option:"mandatory"`
	address           string        `option:"mandatory"`
	readHeaderTimeout time.Duration `option:"mandatory"`
}

func NewServer(opts Options) (*chi.Mux, *http.Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate options error: %v", err)
	}

	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "endpoint not found"}`))
	})

	r.Use(middleware.RequestID)
	//r.Use(customMiddleware.LoggerWithLevel("device"))
	r.Use(logMiddleware.NewWithFilters(
		opts.logger,
		func() logMiddleware.Filter {
			if !opts.isProduction {
				return logMiddleware.IgnoreStatus() //Исключает запросы без ошибок, но при debug: true передает всё
			}
			return logMiddleware.IgnoreStatus(200, 204, 404)
		}(),
		logMiddleware.IgnorePath("/metrics", "/api/v1/auth/login"), //Исключает метрики из логов
	))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Use(customMiddleware.MetricMiddleware)

	r.Mount("/debug", middleware.Profiler()) // pprof
	r.Mount("/metrics", promhttp.Handler())  // подключение метрик

	httpServer := &http.Server{
		Addr:              opts.address,
		Handler:           r,
		ReadHeaderTimeout: opts.readHeaderTimeout,
	}

	return r, httpServer, nil
}
