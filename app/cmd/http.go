package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"log/slog"
	"os"
	"template/config"
	"template/internal/domain/user"
	"template/pkg/fxslog"
	"template/pkg/httpserver"
	"template/pkg/slorage/postgres"

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
			fxslog.NewLogger,
		),

		postgres.Module,

		//DOMAIN - тут происходит подключение доменов
		user.Module,

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
	rootCmd.AddCommand(seedCmd)
}
