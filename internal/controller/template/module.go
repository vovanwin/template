package template

import (
	"context"

	"github.com/vovanwin/template/internal/pkg/server"
	templatepb "github.com/vovanwin/template/pkg/template"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module возвращает fx.Option для подключения TemplateService.
func Module() fx.Option {
	return fx.Options(
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
