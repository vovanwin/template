package dependency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"app/config"
	customMiddleware "app/internal/shared/middleware"
	"app/pkg/httpserver"
	"app/pkg/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/platform/pkg/logger"
	"go.uber.org/fx"
)

func ProvideConfig() (*config.Config, error) {
	return config.NewConfig()
}

func ProvideLogger(config *config.Config) error {
	err := logger.Init(config.Level, false)
	if err != nil {
		return err
	}
	return nil
}

func ProvideServer(lifecycle fx.Lifecycle, config *config.Config) (*chi.Mux, error) {
	// Объявляю нужные мне милдвары для сервера
	middlewareCustom := func(r *chi.Mux) {
		r.Use(middleware.RequestID)
		// r.Use(customMiddleware.LoggerWithLevel("device"))

		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)

		r.Use(customMiddleware.MetricsMiddleware)
		r.Use(customMiddleware.TracingMiddleware)

		// Admin endpoints for runtime log level management
		r.Route(
			"/admin/log", func(r chi.Router) {
				r.Get(
					"/level", func(w http.ResponseWriter, r *http.Request) {
						_ = json.NewEncoder(w).Encode(map[string]string{"level": logger.Level()})
					},
				)
				r.Post(
					"/level", func(w http.ResponseWriter, r *http.Request) {
						var req struct {
							Level string `json:"level"`
						}
						if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Level == "" {
							w.WriteHeader(http.StatusBadRequest)
							_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid level payload"})
							return
						}
						logger.SetLevel(req.Level)
						_ = json.NewEncoder(w).Encode(map[string]string{"level": logger.Level()})
					},
				)
			},
		)

		//chi.Mount("/debug", middleware.Profiler()) // для дебага
	}

	opt := httpserver.NewOptions(
		config.Address(),
		config.ReadHeaderTimeout,
		httpserver.WithMiddlewareSetup(middlewareCustom),
	)
	router, server, err := httpserver.NewServer(opt)
	if err != nil {
		return nil, fmt.Errorf("create http server: %w", err)
	}
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if config.IsLocal() {
					// 👇 выводит все роуты в консоль🚶‍♂️
					httpserver.PrintAllRegisteredRoutes(router)
				}

				go func() {
					lg := logger.Named("http-server")
					lg.Info(context.Background(), "Сервер запущен")
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						lg.Error(context.Background(), "Ошибка запуска сервера")
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				lg := logger.Named("http-server")
				lg.Info(ctx, "Выключение...")
				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
				defer cancel()
				if err := server.Shutdown(shutdownCtx); err != nil {
					lg.Error(ctx, "Ошибка при остановке сервера")
				}
				return nil
			},
		},
	)

	return router, nil
}

func ProvidePgx(config *config.Config) (*postgres.Postgres, error) {
	opt := postgres.NewOptions(
		config.PG.HostPG,
		config.PG.UserPG,
		config.PG.PasswordPG,
		config.PG.DbNamePG,
		config.PG.PortPG,
		config.PG.SchemePG,
		config.IsProduction(),
	)

	connect, err := postgres.New(opt)
	if err != nil {
		return nil, fmt.Errorf("create gorm connection: %w", err)
	}

	return connect, nil
}
