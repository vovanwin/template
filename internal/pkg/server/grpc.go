package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/vovanwin/template/internal/pkg/grpc/health"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) initGRPC(log *slog.Logger) error {
	s.grpcServer = grpc.NewServer(s.grpcOptions...)

	for _, reg := range s.grpcRegistrators {
		reg(s.grpcServer)
	}

	health.RegisterService(s.grpcServer)
	reflection.Register(s.grpcServer)

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.GRPCPort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	go func() {
		log.Info("gRPC сервер запущен", slog.String("addr", addr))
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Error("gRPC сервер остановлен с ошибкой", slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (s *Server) stopGRPC(log *slog.Logger) {
	if s.grpcServer != nil {
		log.Info("gRPC сервер завершает работу...")
		s.grpcServer.GracefulStop()
	}
}
