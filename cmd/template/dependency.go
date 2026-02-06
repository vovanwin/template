package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/pkg/logger"
	postgres2 "github.com/vovanwin/template/internal/pkg/storage/postgres"
)

func ProvideConfig(configDir string) func() (*config.Config, error) {
	return func() (*config.Config, error) {
		cfg, err := config.Load(&config.LoadOptions{
			ConfigDir: configDir,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("=== Конфигурация загружена ===")
		fmt.Printf("Окружение: %s\n", cfg.GetEnv())
		fmt.Printf("Уровень логирования: %s\n", cfg.Log.Level)
		fmt.Println()
		return cfg, nil
	}
}

func ProvideLogger(cfg *config.Config) *slog.Logger {
	logger := logger.NewLogger(logger.Options{
		Level: cfg.Log.Level,
		JSON:  cfg.Log.Format,
	})
	logger.Debug("start logger")
	return logger
}

//// InitLogger инициализирует логгер для fx.Provide
//func InitLogger(config *config.Config) *struct{} {
//	err := logger.Init(config.Level, false)
//	if err != nil {
//		panic(err) // В fx.Provide лучше паниковать при критических ошибках инициализации
//	}
//	return &struct{}{} // Возвращаем пустую структуру
//}

// ProvideJWTService создает JWT сервис
//func ProvideJWTService(config *config.Config) jwt.JWTService {
//	return jwt.NewJWTService(config.JWT.SignKey, config.JWT.TokenTTL, config.JWT.RefreshTTL)
//}

//// ProvideDebugServer запускает отдельный debug/admin HTTP сервер
//func ProvideDebugServer(lifecycle fx.Lifecycle, config *config.Config) error {
//
//	r := chi.NewRouter()
//	r.Use(middleware.Recoverer)
//	r.Use(middleware.RequestID)
//
//	// Профилирование
//	r.Mount("/debug", middleware.Profiler())
//
//	// Admin endpoints for runtime log level management
//	r.Route(
//		"/admin/log", func(r chi.Router) {
//			r.Get(
//				"/level", func(w http.ResponseWriter, r *http.Request) {
//					_ = json.NewEncoder(w).Encode(map[string]string{"level": logger.Level()})
//				},
//			)
//			r.Post(
//				"/level", func(w http.ResponseWriter, r *http.Request) {
//					var req struct {
//						Level string `json:"level"`
//					}
//					if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Level == "" {
//						w.WriteHeader(http.StatusBadRequest)
//						_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid level payload"})
//						return
//					}
//					logger.SetLevel(req.Level)
//					_ = json.NewEncoder(w).Encode(map[string]string{"level": logger.Level()})
//				},
//			)
//		},
//	)
//
//	srv := &http.Server{
//		Addr:    net.JoinHostPort(config.Server.Host, config.Server.DebugPort),
//		Handler: r,
//	}
//
//	lifecycle.Append(
//		fx.Hook{
//			OnStart: func(ctx context.Context) error {
//				go func() {
//					lg := logger.Named("debug-server")
//					lg.Info(context.Background(), fmt.Sprintf("Debug сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.DebugPort)))
//					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//						lg.Error(context.Background(), "Ошибка запуска debug сервера", zap.Error(err))
//					}
//				}()
//				return nil
//			},
//			OnStop: func(ctx context.Context) error {
//				lg := logger.Named("debug-server")
//				lg.Info(context.Background(), "Debug сервер завершает работу...")
//				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
//				defer cancel()
//				return srv.Shutdown(shutdownCtx)
//			},
//		},
//	)
//
//	return nil
//}
//
//// ProvideSwaggerServer запускает сервер со Swagger UI
//func ProvideSwaggerServer(lifecycle fx.Lifecycle, config *config.Config) error {
//	// Проверяем, включен ли Swagger сервер
//
//	r := chi.NewRouter()
//
//	// Раздаём всю директорию со спеками, чтобы $ref ссылки работали
//	specDir := filepath.Join("..", "shared", "api", "app", "v1")
//	fileServer := http.StripPrefix("/spec/", http.FileServer(http.Dir(specDir)))
//	r.Handle("/spec/*", fileServer)
//
//	// Простая страница Swagger UI с CDN, указываем на главный файл в каталоге
//	r.Get(
//		"/", func(w http.ResponseWriter, req *http.Request) {
//			w.Header().Set("Content-Type", "text/html; charset=utf-8")
//
//			_, _ = w.Write(
//				[]byte(`<!doctype html> <!-- Important: must specify -->
//<html>
//  <head>
//    <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 characters -->
//    <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
//  </head>
//  <body>
//    <rapi-doc
//      spec-url = "/spec/app.v1.swagger.yml"
//    > </rapi-doc>
//  </body>
//</html>`),
//			)
//		},
//	)
//
//	srv := &http.Server{
//		Addr:    net.JoinHostPort(config.Server.Host, config.Server.SwaggerPort),
//		Handler: r,
//	}
//
//	lifecycle.Append(
//		fx.Hook{
//			OnStart: func(ctx context.Context) error {
//				go func() {
//					lg := logger.Named("swagger-server")
//					lg.Info(context.Background(), fmt.Sprintf("Swagger сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.SwaggerPort)))
//					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//						lg.Error(context.Background(), "Ошибка запуска swagger сервера", zap.Error(err))
//					}
//				}()
//				return nil
//			},
//			OnStop: func(ctx context.Context) error {
//				lg := logger.Named("swagger-server")
//				lg.Info(context.Background(), "Swagger сервер завершает работу...")
//				shutdownCtx, cancel := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
//				defer cancel()
//				return srv.Shutdown(shutdownCtx)
//			},
//		},
//	)
//
//	return nil
//}

//// ProvideGRPCServer запускает gRPC сервер
//func ProvideGRPCServer(lifecycle fx.Lifecycle, config *config.Config) error {
//
//	lis, err := net.Listen("tcp", net.JoinHostPort(config.Server.Host, config.Server.GRPCPort))
//	if err != nil {
//		return fmt.Errorf("listen grpc: %w", err)
//	}
//	s := grpc.NewServer()
//	healthsvc.RegisterService(s)
//	reflection.Register(s)
//
//	lifecycle.Append(
//		fx.Hook{
//			OnStart: func(ctx context.Context) error {
//				go func() {
//					lg := logger.Named("grpc-server")
//					lg.Info(context.Background(), fmt.Sprintf("gRPC сервер запущен на %s", net.JoinHostPort(config.Server.Host, config.Server.GRPCPort)))
//					if err := s.Serve(lis); err != nil {
//						lg.Error(context.Background(), "Ошибка запуска gRPC сервера", zap.Error(err))
//					}
//				}()
//				return nil
//			},
//			OnStop: func(ctx context.Context) error {
//				lg := logger.Named("grpc-server")
//				lg.Info(context.Background(), "gRPC сервер завершает работу...")
//				done := make(chan struct{})
//				go func() { s.GracefulStop(); close(done) }()
//				select {
//				case <-done:
//					return nil
//				case <-time.After(config.GracefulTimeout * time.Second):
//					lg.Info(context.Background(), "Принудительная остановка gRPC сервера")
//					s.Stop()
//					return nil
//				}
//			},
//		},
//	)
//
//	return nil
//}

func ProvidePgx(c *config.Config) (*postgres2.Postgres, error) {
	opt := postgres2.NewOptions(
		c.PG.Host,
		c.PG.User,
		c.PG.Password,
		c.PG.Db,
		c.PG.Port,
		c.PG.Scheme,
		config.IsProduction(),
	)

	connect, err := postgres2.New(opt)
	if err != nil {
		return nil, fmt.Errorf("create pgx connection: %w", err)
	}

	return connect, nil
}

//
//// ProvideTemporal создает Temporal сервис
//func ProvideTemporal(lifecycle fx.Lifecycle, config *config.Config, _ *struct{}) *temporal.Service {
//	serviceConfig := temporal.ServiceConfig{
//		Client: temporal.Config{
//			Host:      config.Temporal.Host,
//			Port:      config.Temporal.Port,
//			Namespace: config.Temporal.Namespace,
//		},
//		Worker: temporal.WorkerConfig{
//			TaskQueue: config.Temporal.TaskQueue,
//		},
//	}
//
//	service, err := temporal.NewService(serviceConfig)
//	if err != nil {
//		lg := logger.Named("temporal-service")
//		lg.Error(context.Background(), "Ошибка создания Temporal сервиса", zap.Error(err))
//		return nil
//	}
//
//	// Регистрируем lifecycle hooks
//	lifecycle.Append(fx.Hook{
//		OnStart: func(ctx context.Context) error {
//			lg := logger.Named("temporal-service")
//			lg.Info(ctx, fmt.Sprintf("Запуск Temporal сервиса (host: %s:%d, namespace: %s)",
//				config.Temporal.Host, config.Temporal.Port, config.Temporal.Namespace))
//			return service.Start(ctx)
//		},
//		OnStop: func(ctx context.Context) error {
//			lg := logger.Named("temporal-service")
//			lg.Info(ctx, "Temporal сервис завершает работу...")
//			service.Stop(ctx)
//			return nil
//		},
//	})
//
//	return service
//}
