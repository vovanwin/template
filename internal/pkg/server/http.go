package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func (s *Server) initHTTP(log *slog.Logger) error {
	gwMux := runtime.NewServeMux()

	for _, reg := range s.gatewayRegistrators {
		if err := reg(context.Background(), gwMux, s.grpcServer); err != nil {
			return fmt.Errorf("register gateway: %w", err)
		}
	}

	r := chi.NewRouter()

	// Логирование запросов через slog
	r.Use(SlogRequestLogger(log))

	// Пользовательские middleware
	for _, mw := range s.httpMiddleware {
		r.Use(mw)
	}

	r.Mount("/", gwMux)

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.HTTPPort)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info("HTTP gateway запущен", slog.String("addr", addr))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP gateway остановлен с ошибкой", slog.String("error", err.Error()))
		}
	}()

	return nil
}

func (s *Server) stopHTTP(ctx context.Context, log *slog.Logger) error {
	if s.httpServer != nil {
		log.Info("HTTP gateway завершает работу...")
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
