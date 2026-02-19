package auth

import (
	"context"

	"github.com/vovanwin/platform/server"
	authpb "github.com/vovanwin/template/pkg/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module возвращает fx.Option для подключения AuthService.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewAuthGRPCServer),
		fx.Provide(
			fx.Annotate(
				func(srv *AuthGRPCServer) server.GRPCRegistrator {
					return func(s *grpc.Server) {
						authpb.RegisterAuthServiceServer(s, srv)
					}
				},
				fx.ResultTags(`group:"grpc_registrators"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(srv *AuthGRPCServer) server.GatewayRegistrator {
					return func(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
						return authpb.RegisterAuthServiceHandlerServer(ctx, mux, srv)
					}
				},
				fx.ResultTags(`group:"gateway_registrators"`),
			),
		),
	)
}
