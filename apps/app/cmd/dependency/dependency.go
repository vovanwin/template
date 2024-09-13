package dependency

import (
	"app/config"
	customMiddleware "app/internal/shared/middleware"
	"app/pkg/fxslog"
	"app/pkg/httpserver"
	"app/pkg/storage/postgres"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
	"log"
	"log/slog"
	"time"
)

func ProvideConfig() (*config.Config, error) {
	return config.NewConfig()
}

func ProvideLogger(config *config.Config) (*slog.Logger, error) {
	opt := fxslog.NewOptions(fxslog.WithEnv(config.Env), fxslog.WithLevel(config.Level))
	return fxslog.NewLogger(opt)
}

func ProvideServer(lifecycle fx.Lifecycle, logger *slog.Logger, config *config.Config) (*chi.Mux, error) {
	// Объявляю нужные мне милдвары для сервера
	middlewareCustom := func(chi *chi.Mux) {
		chi.Use(middleware.RequestID)
		// r.Use(customMiddleware.LoggerWithLevel("device"))

		chi.Use(middleware.Recoverer)
		chi.Use(middleware.URLFormat)

		chi.Use(customMiddleware.MetricsMiddleware)
		chi.Use(customMiddleware.TracingMiddleware)

		chi.Mount("/debug", middleware.Profiler()) // для дебага
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
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if config.IsLocal() {
				// 👇 выводит все роуты в консоль🚶‍♂️
				httpserver.PrintAllRegisteredRoutes(router)
			}

			go func() {
				log.Printf("Сервер запущен на %s\n", config.Address())
				if err := server.ListenAndServe(); err != nil {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Выключение...")
			shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Println(err)
			}
			return nil
		},
	})

	return router, nil
}

func ProvidePgx(config *config.Config, logger *slog.Logger) (*postgres.Postgres, error) {
	opt := postgres.NewOptions(
		logger,
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
