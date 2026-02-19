package template

import (
	"context"

	templatepb "github.com/vovanwin/template/pkg/template"
)

func (s *TemplateGRPCServer) GetHealth(_ context.Context, req *templatepb.GetHealthRequest) (*templatepb.GetHealthResponse, error) {
	return &templatepb.GetHealthResponse{Status: "ok"}, nil

}
