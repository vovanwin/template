package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vovanwin/platform/logger"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/api"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/pkg/metrics"
	appOtel "github.com/vovanwin/template/internal/pkg/otel"
	postgres2 "github.com/vovanwin/template/internal/pkg/storage/postgres"
	pkg "github.com/vovanwin/template/pkg"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

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

func ProvideOtel(cfg *config.Config, log *slog.Logger) (*appOtel.Provider, error) {
	provider := &appOtel.Provider{}

	if !cfg.Features.EnableTracing && !cfg.Features.EnableMetrics {
		log.Info("OTEL отключён (features.enable_tracing=false, features.enable_metrics=false)")
		return provider, nil
	}

	endpoint := cfg.Otel.Endpoint
	serviceName := cfg.App.Name

	if cfg.Features.EnableTracing {
		tp, err := appOtel.InitTracer(context.Background(), serviceName, endpoint)
		if err != nil {
			return nil, fmt.Errorf("init otel tracer: %w", err)
		}
		provider.TracerProvider = tp
		log.Info("OTEL трейсинг включён", slog.String("endpoint", endpoint))
	}

	if cfg.Features.EnableMetrics {
		mp, err := appOtel.InitMeter(context.Background(), serviceName, endpoint)
		if err != nil {
			return nil, fmt.Errorf("init otel meter: %w", err)
		}
		provider.MeterProvider = mp
		log.Info("OTEL метрики включены", slog.String("endpoint", endpoint))
	}

	return provider, nil
}

func ProvideServerModule(cfg *config.Config) fx.Option {
	opts := []server.Option{
		server.WithHTTPMiddleware(middleware.Recoverer, middleware.RequestID),
	}

	if cfg.Features.EnableMetrics {
		opts = append(opts, server.WithDebugHandler("/metrics", metrics.Handler()))
	}

	if cfg.Features.EnableMetrics || cfg.Features.EnableTracing {
		opts = append(opts,
			server.WithHTTPMiddleware(appOtel.HTTPMiddleware(cfg.App.Name+"-http")),
			server.WithGRPCOptions(grpc.StatsHandler(otelgrpc.NewServerHandler())),
		)
	}

	return server.NewModule(opts...)
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
