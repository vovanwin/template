package httpserver

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"

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

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//logrus.Print("Start Http Server.")
			go func() {

				// 👇 выводит все роуты в консоль🚶‍♂️
				chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
					fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
					return nil
				})

				err := http.ListenAndServe(config.Server.Address, r)
				if err != nil {
					log.Fatal(err)
				}
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
