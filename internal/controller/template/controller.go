package template

import (
	templatepb "github.com/vovanwin/template/pkg/template"
)

// TemplateGRPCServer реализует gRPC сервис TemplateService.
type TemplateGRPCServer struct {
	templatepb.UnimplementedTemplateServiceServer
}
