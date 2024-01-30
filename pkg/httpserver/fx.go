package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"time"

	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"template/config"
)

var Module = fx.Provide(newModule)

func newModule(lifecycle fx.Lifecycle, config config.Config, logger *slog.Logger) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//r.Mount("/debug", middleware.Profiler()) // для дебага
	r.Mount("/metrics", promhttp.Handler()) // подключение метрик

	httpServer := &http.Server{
		Addr:              config.Server.Address,
		Handler:           r,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//logrus.Print("Start Http Server.")
			go func() {

				// 👇 выводит все роуты в консоль🚶‍♂️
				printAllRegisteredRoutes(r)

				go func() {
					start(httpServer, config)
				}()

			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			//logrus.Print("Stopping Http Server.")
			return nil
		},
	})

	return r
}

func start(s *http.Server, config config.Config) {
	log.Printf("Сервер запущен на %s\n", config.Server.Address)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, config config.Config, s *http.Server) error {
	log.Println("Выключение...")

	ctx, shutdown := context.WithTimeout(ctx, config.Server.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}

	return nil
}
