package main

import (
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/platform/logger"
	platformotel "github.com/vovanwin/platform/otel"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/api"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/pkg/storage/postgres"
	"github.com/vovanwin/template/pkg"

	"go.uber.org/fx"
)

func ProvideLogger(cfg *config.Config, lc fx.Lifecycle) (*slog.Logger, error) {
	l, closer := logger.NewLogger(logger.Options{
		Level:       cfg.Log.Level,
		JSON:        cfg.Log.Format,
		LokiEnabled: cfg.Log.LokiEnabled,
		LokiURL:     cfg.Log.LokiUrl,
		ServiceName: cfg.App.Name,
		TraceID:     cfg.Metrics.EnableMetrics,
	})
	lc.Append(fx.StopHook(closer))
	l.Debug("start logger")
	return l, nil
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

func ProvideServerModule(cfg *config.Config) fx.Option {
	opts := []server.Option{
		server.WithHTTPMiddleware(middleware.RequestID),
	}

	if cfg.Metrics.EnableMetrics {
		opts = append(opts, server.WithOtel(platformotel.Config{
			ServiceName: cfg.App.Name,
			Endpoint:    cfg.Otel.Endpoint,
			SampleRate:  cfg.Otel.SampleRate,
		}))
	}

	return server.NewModule(opts...)
}

func ProvidePgx(c *config.Config) (*postgres.Postgres, error) {
	opt := postgres.NewOptions(
		c.PG.Host,
		c.PG.User,
		c.PG.Password,
		c.PG.Db,
		c.PG.Port,
		c.PG.Scheme,
		config.IsProduction(),
	)

	connect, err := postgres.New(opt)
	if err != nil {
		return nil, fmt.Errorf("create pgx connection: %w", err)
	}

	return connect, nil
}
