package dependency

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vovanwin/template/app/config"
	customMiddleware "github.com/vovanwin/template/app/internal/shared/middleware"
	"github.com/vovanwin/template/app/pkg/httpserver"
	"github.com/vovanwin/template/app/pkg/jwt"
	"github.com/vovanwin/template/app/pkg/storage/postgres"

	"github.com/vovanwin/platform/pkg/temporal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	healthsvc "github.com/vovanwin/platform/pkg/grpc/health"
	"github.com/vovanwin/platform/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

// InitLogger инициализирует логгер для fx.Provide
func InitLogger(config *config.Config) *struct{} {
	err := logger.Init(config.Level, false)
	if err != nil {
		panic(err) // В fx.Provide лучше паниковать при критических ошибках инициализации
	}
	return &struct{}{} // Возвращаем пустую структуру
}

// ProvideJWTService создает JWT сервис
func ProvideJWTService(config *config.Config) jwt.JWTService {
	return jwt.NewJWTService(config.JWT.SignKey, config.JWT.TokenTTL, config.JWT.RefreshTTL)
}

func ProvideServer(lifecycle fx.Lifecycle, config *config.Config) (*chi.Mux, error) {
	// Объявляю нужные мне милдвары для сервера
	// Создаем rate limiter
	rateLimiter := customMiddleware.NewRateLimiter()

	middlewareCustom := func(r *chi.Mux) {
		r.Use(middleware.RequestID)

		// CORS для Swagger UI и других клиентов
		r.Use(
			cors.Handler(
				cors.Options{
					AllowedOrigins:   []string{"*"},
					AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
					ExposedHeaders:   []string{},
					AllowCredentials: false,
					MaxAge:           300, // 5 minutes
				},
			),
		)

		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)

		// Rate limiting (раньше других для быстрого отклонения)
		r.Use(rateLimiter.RateLimitMiddleware())

		// JWT защита обрабатывается в ogen SecurityHandler

		r.Use(customMiddleware.MetricsMiddleware)
		r.Use(customMiddleware.TracingMiddleware)
	}

	opt := httpserver.NewOptions(
		net.JoinHostPort(config.Server.Host, "8080"),
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

// ProvideDebugServer запускает отдельный debug/admin HTTP сервер на 8082
func ProvideDebugServer(lifecycle fx.Lifecycle, config *config.Config) error {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Профилирование
	r.Mount("/debug", middleware.Profiler())

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

	srv := &http.Server{
		Addr:    net.JoinHostPort(config.Server.Host, "8082"),
		Handler: r,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					lg := logger.Named("debug-server")
					lg.Info(context.Background(), "Debug server started")
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						lg.Error(context.Background(), "Ошибка запуска debug сервера")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
				defer cancel()
				return srv.Shutdown(shutdownCtx)
			},
		},
	)

	return nil
}

// ProvideSwaggerServer запускает сервер со Swagger UI на 8084
func ProvideSwaggerServer(lifecycle fx.Lifecycle, config *config.Config) error {
	r := chi.NewRouter()

	// Раздаём всю директорию со спеками, чтобы $ref ссылки работали
	specDir := filepath.Join("..", "shared", "api", "app", "v1")
	fileServer := http.StripPrefix("/spec/", http.FileServer(http.Dir(specDir)))
	r.Handle("/spec/*", fileServer)

	// Простая страница Swagger UI с CDN, указываем на главный файл в каталоге
	r.Get(
		"/", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			_, _ = w.Write(
				[]byte(`<!doctype html> <!-- Important: must specify -->
<html>
  <head>
    <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 characters -->
    <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
  </head>
  <body>
    <rapi-doc
      spec-url = "/spec/app.v1.swagger.yml"
    > </rapi-doc>
  </body>
</html>`),
			)
		},
	)

	srv := &http.Server{
		Addr:    net.JoinHostPort(config.Server.Host, "8084"),
		Handler: r,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					lg := logger.Named("swagger-server")
					lg.Info(context.Background(), "Swagger server started")
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						lg.Error(context.Background(), "Ошибка запуска swagger сервера")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
				defer cancel()
				return srv.Shutdown(shutdownCtx)
			},
		},
	)

	return nil
}

// ProvideGRPCServer запускает gRPC сервер на 8081
func ProvideGRPCServer(lifecycle fx.Lifecycle, config *config.Config) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(config.Server.Host, "8081"))
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}
	s := grpc.NewServer()
	healthsvc.RegisterService(s)
	reflection.Register(s)

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					lg := logger.Named("grpc-server")
					lg.Info(context.Background(), "gRPC server started")
					if err := s.Serve(lis); err != nil {
						lg.Error(context.Background(), "Ошибка запуска grpc сервера")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				done := make(chan struct{})
				go func() { s.GracefulStop(); close(done) }()
				select {
				case <-done:
					return nil
				case <-time.After(config.GracefulTimeout * time.Second):
					s.Stop()
					return nil
				}
			},
		},
	)

	return nil
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
		return nil, fmt.Errorf("create pgx connection: %w", err)
	}

	return connect, nil
}

// ProvidePool предоставляет *pgxpool.Pool для sqlc и других компонентов
func ProvidePool(pg *postgres.Postgres) *pgxpool.Pool {
	pool, ok := pg.Pool.(*pgxpool.Pool)
	if !ok {
		panic("unable to cast pool to *pgxpool.Pool")
	}
	return pool
}

// ProvideTemporal создает Temporal сервис
func ProvideTemporal(lifecycle fx.Lifecycle, config *config.Config, _ *struct{}) (*temporal.Service, error) {
	serviceConfig := temporal.ServiceConfig{
		Client: temporal.Config{
			Host:      config.Temporal.Host,
			Port:      config.Temporal.Port,
			Namespace: config.Temporal.Namespace,
		},
		Worker: temporal.WorkerConfig{
			TaskQueue: config.Temporal.TaskQueue,
		},
	}

	service, err := temporal.NewService(serviceConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal service: %w", err)
	}

	// Регистрируем lifecycle hooks
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return service.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			service.Stop(ctx)
			return nil
		},
	})

	return service, nil
}
