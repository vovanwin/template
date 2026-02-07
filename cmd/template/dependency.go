package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/platform/logger"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/api"
	"github.com/vovanwin/template/config"
	postgres2 "github.com/vovanwin/template/internal/pkg/storage/postgres"
	pkg "github.com/vovanwin/template/pkg"

	"go.uber.org/fx"
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
	l := logger.NewLogger(logger.Options{
		Level: cfg.Log.Level,
		JSON:  cfg.Log.Format,
	})
	l.Debug("start logger")
	return l
}

func ProvideServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Host:        cfg.Server.Host,
		GRPCPort:    cfg.Server.GrpcPort,
		HTTPPort:    cfg.Server.HttpPort,
		SwaggerPort: cfg.Server.SwaggerPort,
		DebugPort:   cfg.Server.DebugPort,
		SwaggerFS:   pkg.EmbedSwagger,
		ProtoFS:     api.EmbedProto,
	}
}

func ProvideServerModule() fx.Option {
	return server.NewModule(
		server.WithHTTPMiddleware(middleware.Recoverer, middleware.RequestID),
	)
}

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
