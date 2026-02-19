package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/platform/logger"
	platformotel "github.com/vovanwin/platform/otel"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/api"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/pkg/etcdstore"
	"github.com/vovanwin/template/internal/pkg/flagsui"
	"github.com/vovanwin/template/internal/pkg/jwt"
	authmw "github.com/vovanwin/template/internal/pkg/middleware"
	"github.com/vovanwin/template/internal/pkg/storage/postgres"
	"github.com/vovanwin/template/pkg"

	"go.uber.org/fx"
	"google.golang.org/grpc"
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

// ProvideFlags создает Flags с EtcdStore (если настроен) или MemoryStore (fallback).
// Вызывается до fx.New(), поэтому не зависит от DI-контейнера.
func ProvideFlags(cfg *config.Config) (*config.Flags, func()) {
	if cfg.Etcd.Endpoints != "" {
		endpoints := strings.Split(cfg.Etcd.Endpoints, ",")
		store, err := etcdstore.New(etcdstore.Config{
			Endpoints: endpoints,
			Prefix:    cfg.Etcd.Prefix,
		}, slog.Default())

		if err != nil {
			slog.Warn("etcd unavailable, falling back to MemoryStore", slog.Any("error", err))
			return config.NewFlags(config.NewMemoryStore(config.DefaultFlagValues())), func() {}
		}

		store.SetDefaults(config.DefaultFlagValues())
		_ = store.Watch(context.Background(), func(key string) {
			slog.Info("feature flag changed", slog.String("key", key))
		})

		slog.Info("feature flags: etcd store", slog.String("endpoints", cfg.Etcd.Endpoints))
		return config.NewFlags(store), func() { store.Close() }
	}

	slog.Info("feature flags: memory store")
	return config.NewFlags(config.NewMemoryStore(config.DefaultFlagValues())), func() {}
}

func ProvideServerModule(cfg *config.Config, flags *config.Flags, jwtService jwt.JWTService) fx.Option {
	opts := []server.Option{
		server.WithHTTPMiddleware(middleware.RequestID),
		server.WithDebugHandler("/flags", flagsui.Handler(flags)),
		server.WithDebugHandler("/flags/", flagsui.Handler(flags)),
		server.WithGRPCOptions(grpc.ChainUnaryInterceptor(authmw.AuthInterceptor(jwtService, cfg.Server.AuthBypass))),
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
