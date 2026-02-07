package server

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Config — конфигурация серверов, передаётся потребителем.
type Config struct {
	Host        string
	GRPCPort    string
	HTTPPort    string
	SwaggerPort string
	DebugPort   string
	SwaggerFS   fs.FS // встроенная FS с файлами *.swagger.json
	ProtoFS     fs.FS // встроенная FS с файлами *.proto
}

// GRPCRegistrator — колбэк для регистрации gRPC сервисов.
type GRPCRegistrator func(s *grpc.Server)

// GatewayRegistrator — колбэк для регистрации grpc-gateway in-process.
type GatewayRegistrator func(ctx context.Context, mux *runtime.ServeMux, server *grpc.Server) error

// Option — функциональные опции для Server.
type Option func(*Server)

// Server объединяет все 4 сервера: gRPC, HTTP gateway, Swagger, Debug.
type Server struct {
	cfg Config

	grpcRegistrators    []GRPCRegistrator
	gatewayRegistrators []GatewayRegistrator
	httpMiddleware      []func(http.Handler) http.Handler
	debugMiddleware     []func(http.Handler) http.Handler
	grpcOptions         []grpc.ServerOption

	grpcServer *grpc.Server
	httpServer *http.Server
	swaggerSrv *http.Server
	debugSrv   *http.Server
}

// WithGRPCRegistrator добавляет колбэк для регистрации gRPC сервисов.
func WithGRPCRegistrator(r GRPCRegistrator) Option {
	return func(s *Server) {
		s.grpcRegistrators = append(s.grpcRegistrators, r)
	}
}

// WithGatewayRegistrator добавляет колбэк для регистрации grpc-gateway хендлеров.
func WithGatewayRegistrator(r GatewayRegistrator) Option {
	return func(s *Server) {
		s.gatewayRegistrators = append(s.gatewayRegistrators, r)
	}
}

// WithHTTPMiddleware добавляет middleware на HTTP gateway.
func WithHTTPMiddleware(mw ...func(http.Handler) http.Handler) Option {
	return func(s *Server) {
		s.httpMiddleware = append(s.httpMiddleware, mw...)
	}
}

// WithDebugMiddleware добавляет middleware на Debug сервер.
func WithDebugMiddleware(mw ...func(http.Handler) http.Handler) Option {
	return func(s *Server) {
		s.debugMiddleware = append(s.debugMiddleware, mw...)
	}
}

// WithGRPCOptions добавляет опции для gRPC сервера.
func WithGRPCOptions(opts ...grpc.ServerOption) Option {
	return func(s *Server) {
		s.grpcOptions = append(s.grpcOptions, opts...)
	}
}

func newServer(cfg Config, opts ...Option) *Server {
	s := &Server{cfg: cfg}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
