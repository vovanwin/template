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
	"go.uber.org/zap"
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
	// Проверяем, включен ли HTTP сервер
	if !config.Server.EnableHTTP {
		lg := logger.Named("http-server")
		lg.Info(context.Background(), "HTTP сервер отключен конфигурацией")
		return chi.NewRouter(), nil // Возвращаем пустой роутер
	}

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

	}

	opt := httpserver.NewOptions(
		net.JoinHostPort(config.Server.Host, config.Server.HTTPPort),
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
					lg.Info(context.Background(), fmt.Sprintf("HTTP сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.HTTPPort)))
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						lg.Error(context.Background(), "Ошибка запуска HTTP сервера", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				lg := logger.Named("http-server")
				lg.Info(ctx, "HTTP сервер завершает работу...")
				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
				defer cancel()
				if err := server.Shutdown(shutdownCtx); err != nil {
					lg.Error(ctx, "Ошибка остановки HTTP сервера", zap.Error(err))
				}
				return nil
			},
		},
	)

	return router, nil
}

// ProvideDebugServer запускает отдельный debug/admin HTTP сервер
func ProvideDebugServer(lifecycle fx.Lifecycle, config *config.Config) error {
	// Проверяем, включен ли Debug сервер
	if !config.Server.EnableDebug {
		lg := logger.Named("debug-server")
		lg.Info(context.Background(), "Debug сервер отключен конфигурацией")
		return nil
	}

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
		Addr:    net.JoinHostPort(config.Server.Host, config.Server.DebugPort),
		Handler: r,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					lg := logger.Named("debug-server")
					lg.Info(context.Background(), fmt.Sprintf("Debug сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.DebugPort)))
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						lg.Error(context.Background(), "Ошибка запуска debug сервера", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				lg := logger.Named("debug-server")
				lg.Info(context.Background(), "Debug сервер завершает работу...")
				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
				defer cancel()
				return srv.Shutdown(shutdownCtx)
			},
		},
	)

	return nil
}

// ProvideSwaggerServer запускает сервер со Swagger UI
func ProvideSwaggerServer(lifecycle fx.Lifecycle, config *config.Config) error {
	// Проверяем, включен ли Swagger сервер
	if !config.Server.EnableSwagger {
		lg := logger.Named("swagger-server")
		lg.Info(context.Background(), "Swagger сервер отключен конфигурацией")
		return nil
	}

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
		Addr:    net.JoinHostPort(config.Server.Host, config.Server.SwaggerPort),
		Handler: r,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					lg := logger.Named("swagger-server")
					lg.Info(context.Background(), fmt.Sprintf("Swagger сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.SwaggerPort)))
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						lg.Error(context.Background(), "Ошибка запуска swagger сервера", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				lg := logger.Named("swagger-server")
				lg.Info(context.Background(), "Swagger сервер завершает работу...")
				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
				defer cancel()
				return srv.Shutdown(shutdownCtx)
			},
		},
	)

	return nil
}

// ProvideGRPCServer запускает gRPC сервер
func ProvideGRPCServer(lifecycle fx.Lifecycle, config *config.Config) error {
	// Проверяем, включен ли gRPC сервер
	if !config.Server.EnableGRPC {
		lg := logger.Named("grpc-server")
		lg.Info(context.Background(), "gRPC сервер отключен конфигурацией")
		return nil
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(config.Server.Host, config.Server.GRPCPort))
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
					lg.Info(context.Background(), fmt.Sprintf("gRPC сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.GRPCPort)))
					if err := s.Serve(lis); err != nil {
						lg.Error(context.Background(), "Ошибка запуска gRPC сервера", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				lg := logger.Named("grpc-server")
				lg.Info(context.Background(), "gRPC сервер завершает работу...")
				done := make(chan struct{})
				go func() { s.GracefulStop(); close(done) }()
				select {
				case <-done:
					return nil
				case <-time.After(config.GracefulTimeout * time.Second):
					lg.Info(context.Background(), "Принудительная остановка gRPC сервера")
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
func ProvideTemporal(lifecycle fx.Lifecycle, config *config.Config, _ *struct{}) *temporal.Service {
	// Проверяем, включен ли Temporal
	if !config.Server.EnableTemporal {
		lg := logger.Named("temporal-service")
		lg.Info(context.Background(), "Temporal сервис отключен конфигурацией")
		return nil // Возвращаем nil, но без ошибки
	}

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
		lg := logger.Named("temporal-service")
		lg.Error(context.Background(), "Ошибка создания Temporal сервиса", zap.Error(err))
		return nil
	}

	// Регистрируем lifecycle hooks
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lg := logger.Named("temporal-service")
			lg.Info(ctx, fmt.Sprintf("Запуск Temporal сервиса (host: %s:%d, namespace: %s)",
				config.Temporal.Host, config.Temporal.Port, config.Temporal.Namespace))
			return service.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			lg := logger.Named("temporal-service")
			lg.Info(ctx, "Temporal сервис завершает работу...")
			service.Stop(ctx)
			return nil
		},
	})

	return service
}
