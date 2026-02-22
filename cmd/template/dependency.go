package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/platform/logger"
	platformotel "github.com/vovanwin/platform/otel"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/api"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/pkg/centrifugo"
	"github.com/vovanwin/template/internal/pkg/etcdstore"
	"github.com/vovanwin/template/internal/pkg/events"
	"github.com/vovanwin/template/internal/pkg/flagsui"
	"github.com/vovanwin/template/internal/pkg/jwt"
	"github.com/vovanwin/template/internal/pkg/logx"
	authmw "github.com/vovanwin/template/internal/pkg/middleware"
	"github.com/vovanwin/template/internal/pkg/storage/postgres"
	"github.com/vovanwin/template/internal/pkg/temporal"
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

	// Оборачиваем логгер в ComponentHandler для поддержки переопределений уровней
	handler := logx.NewComponentHandler(l.Handler(), cfg.Log.Overrides)
	l = slog.New(handler)

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
	csrfMiddleware := func(next http.Handler) http.Handler {
		return next
	} // Временно отключаем CSRF для теста

	opts := []server.Option{
		server.WithHTTPMiddleware(middleware.RequestID),
		server.WithHTTPMiddleware(authmw.CookieAuthMiddleware),
		server.WithHTTPMiddleware(csrfMiddleware),
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

func ProvideCentrifugoClient(cfg *config.Config, log *slog.Logger) *centrifugo.Client {
	return centrifugo.NewClient(
		cfg.Centrifugo.Addr,
		cfg.Centrifugo.ApiKey,
		cfg.Centrifugo.TokenSecret,
		cfg.Centrifugo.TokenTtl,
		log.With("component", "centrifugo"),
	)
}

func ProvideEventBus(client *centrifugo.Client, log *slog.Logger) *events.Bus {
	return events.NewBus(client, log.With("component", "events"))
}

func ProvideCentrifugoURL(cfg *config.Config) string {
	return cfg.Centrifugo.Addr
}

func ProvideTemporalService(cfg *config.Config, log *slog.Logger) (*temporal.Service, error) {
	return temporal.NewService(temporal.ServiceConfig{
		Client: temporal.Config{
			Host:      cfg.Temporal.Host,
			Port:      cfg.Temporal.Port,
			Namespace: cfg.Temporal.Namespace,
		},
		Worker: temporal.WorkerConfig{TaskQueue: cfg.Temporal.TaskQueue},
	}, log.With("component", "temporal"))
}

func ProvidePgx(c *config.Config, log *slog.Logger) (*postgres.Postgres, error) {
	opt := postgres.NewOptions(
		c.PG.Host,
		c.PG.User,
		c.PG.Password,
		c.PG.Db,
		c.PG.Port,
		c.PG.Scheme,
		config.IsProduction(),
	)

	connect, err := postgres.New(opt, log.With("component", "db"))
	if err != nil {
		return nil, fmt.Errorf("create pgx connection: %w", err)
	}

	return connect, nil
}
