package server

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

// NewModule создаёт fx.Module для серверного пакета.
// Потребитель должен предоставить server.Config и *slog.Logger через fx.Provide.
func NewModule(opts ...Option) fx.Option {
	return fx.Module("server",
		fx.Invoke(func(lc fx.Lifecycle, cfg Config, log *slog.Logger) {
			s := newServer(cfg, opts...)

			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					if err := s.initGRPC(log); err != nil {
						return err
					}
					if err := s.initHTTP(log); err != nil {
						return err
					}
					s.initSwagger(log)
					s.initDebug(log)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					_ = s.stopHTTP(ctx, log)
					_ = s.stopSwagger(ctx, log)
					_ = s.stopDebug(ctx, log)
					s.stopGRPC(log)
					return nil
				},
			})
		}),
	)
}
