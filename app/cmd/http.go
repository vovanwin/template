package cmd

import (
	"app/cmd/dependency"
	"app/cmd/migrateCmd"
	"app/config"
	"app/internal/module/healthcheck"
	"app/internal/module/users"
	"app/internal/shared/middleware"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"log/slog"
	"os"
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
		//fx.NopLogger,

		fx.Provide(
			dependency.ProvideConfig,
			dependency.ProvideLogger,
			dependency.ProvideServer,
		),

		users.Module,

		//  healthcheck
		fx.Invoke(healthcheck.Controller),
		// загружаю мидлваре в приложение
		fx.Provide(middleware.NewMiddleware),

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

func init() {
	rootCmd.AddCommand(migrateCmd.MigrationsCmd)
}
