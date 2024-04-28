package cmd

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/module/healthcheck"
	"github.com/vovanwin/template/pkg/fxslog"
	"github.com/vovanwin/template/pkg/httpserver"
	"go.uber.org/fx"
	"log"
	"log/slog"
	"os"
	"time"
)

var (
	Version = "0.1"

	rootCmd = &cobra.Command{
		Use:     "server",
		Version: Version,
		Short:   "Запуск Http REST API",
		Run: func(cmd *cobra.Command, args []string) {
			fx.New(inject()).Run()
		},
	}
)

func inject() fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewConfig,
			provideLogger,
			provideServer, // TODO: из -за особенностей fx нужно вызвать какой либо контроллер например fx.Invoke(healthcheck.Controller) чтобы выполнилнилась иницыализация сервера
		),

		//  healthcheck
		fx.Invoke(healthcheck.Controller),

		fx.Decorate(func(logger *slog.Logger, config *config.Config) *slog.Logger {
			return logger.
				With("environment", config.Env).
				With("release", Version)
		}),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func provideLogger(config *config.Config) (*slog.Logger, error) {
	opt := fxslog.NewOptions(fxslog.WithEnv(config.Env), fxslog.WithLevel(config.Level))
	return fxslog.NewLogger(opt)
}

func provideServer(lifecycle fx.Lifecycle, logger *slog.Logger, config *config.Config) (*chi.Mux, error) {
	opt := httpserver.NewOptions(logger, config.IsProduction(), config.Address(), config.ReadHeaderTimeout)
	router, server, err := httpserver.NewServer(opt)
	if err != nil {
		return nil, fmt.Errorf("create http server: %v", err)
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			if !config.IsProduction() {
				// 👇 выводит все роуты в консоль🚶‍♂️
				httpserver.PrintAllRegisteredRoutes(router)
			}

			go func() {
				log.Printf("Сервер запущен на %s\n", config.Address())
				err := server.ListenAndServe()
				if err != nil {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Выключение...")

			ctx, shutdown := context.WithTimeout(ctx, config.GracefulTimeout*time.Second)
			defer shutdown()

			err := server.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}

			return nil
		},
	})

	return router, nil
}
