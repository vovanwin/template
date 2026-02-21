package template

import (
	"context"
	"log/slog"

	"github.com/vovanwin/platform/server"
	templatepb "github.com/vovanwin/template/pkg/template"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module возвращает fx.Option для подключения TemplateService.
func Module() fx.Option {
	return fx.Module("api:template",
		fx.Decorate(func(log *slog.Logger) *slog.Logger {
			return log.With("component", "api")
		}),
		fx.Provide(NewTemplateGRPCServer),
		fx.Provide(
			fx.Annotate(
				func(srv *TemplateGRPCServer) server.GRPCRegistrator {
					return func(s *grpc.Server) {
						templatepb.RegisterTemplateServiceServer(s, srv)
					}
				},
				fx.ResultTags(`group:"grpc_registrators"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(srv *TemplateGRPCServer) server.GatewayRegistrator {
					return func(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
						return templatepb.RegisterTemplateServiceHandlerServer(ctx, mux, srv)
					}
				},
				fx.ResultTags(`group:"gateway_registrators"`),
			),
		),
	)
}
