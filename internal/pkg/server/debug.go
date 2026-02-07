package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func (s *Server) initDebug(log *slog.Logger) {
	r := chi.NewRouter()

	// Логирование запросов через slog
	r.Use(SlogRequestLogger(log))

	// Пользовательские middleware
	for _, mw := range s.debugMiddleware {
		r.Use(mw)
	}

	// pprof
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Liveness/readiness для k8s
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.DebugPort)

	s.debugSrv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info("Debug сервер запущен", slog.String("addr", addr))
		if err := s.debugSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Debug сервер остановлен с ошибкой", slog.String("error", err.Error()))
		}
	}()
}

func (s *Server) stopDebug(ctx context.Context, log *slog.Logger) error {
	if s.debugSrv != nil {
		log.Info("Debug сервер завершает работу...")
		return s.debugSrv.Shutdown(ctx)
	}
	return nil
}
