package server

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

// moduleParams — зависимости серверного модуля, включая группы регистраторов.
type moduleParams struct {
	fx.In

	LC                  fx.Lifecycle
	Cfg                 Config
	Log                 *slog.Logger
	GRPCRegistrators    []GRPCRegistrator    `group:"grpc_registrators"`
	GatewayRegistrators []GatewayRegistrator `group:"gateway_registrators"`
}

// NewModule создаёт fx.Module для серверного пакета.
// gRPC и gateway регистраторы собираются автоматически через fx groups.
// Потребитель должен предоставить server.Config и *slog.Logger через fx.Provide.
func NewModule(opts ...Option) fx.Option {
	return fx.Module("server",
		fx.Invoke(func(p moduleParams) {
			s := newServer(p.Cfg, opts...)

			s.grpcRegistrators = append(s.grpcRegistrators, p.GRPCRegistrators...)
			s.gatewayRegistrators = append(s.gatewayRegistrators, p.GatewayRegistrators...)

			p.LC.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					if err := s.initGRPC(p.Log); err != nil {
						return err
					}
					if err := s.initHTTP(p.Log); err != nil {
						return err
					}
					s.initSwagger(p.Log)
					s.initDebug(p.Log)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					_ = s.stopHTTP(ctx, p.Log)
					_ = s.stopSwagger(ctx, p.Log)
					_ = s.stopDebug(ctx, p.Log)
					s.stopGRPC(p.Log)
					return nil
				},
			})
		}),
	)
}
