package template

import (
	"log/slog"

	templatepb "github.com/vovanwin/template/pkg/template"
	"go.uber.org/fx"
)

// Deps содержит зависимости для TemplateGRPCServer.
type Deps struct {
	fx.In

	Log *slog.Logger
}

// TemplateGRPCServer реализует gRPC сервис TemplateService.
type TemplateGRPCServer struct {
	templatepb.UnimplementedTemplateServiceServer
	log *slog.Logger
}

// NewTemplateGRPCServer создаёт новый TemplateGRPCServer.
func NewTemplateGRPCServer(deps Deps) *TemplateGRPCServer {
	return &TemplateGRPCServer{log: deps.Log}
}
