package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
	"template/config"
	"template/internal/controller"
	"template/internal/repository"
	"template/internal/service"
	"template/pkg/fxslog"
	"template/pkg/httpserver"
	"template/pkg/postgres"
	"template/pkg/utils"
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
			utils.NewTimeoutContext,
			fxslog.SetupLogger(),
		),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return fxslog.New(logger)
		}),
		postgres.Module,
		repository.Module,
		service.Module,
		controller.Module,
		httpserver.Module,

		fx.Decorate(func(logger *slog.Logger, config config.Config) *slog.Logger {
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

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.AddCommand(testCmd)
}
