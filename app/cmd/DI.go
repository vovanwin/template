package cmd

import (
	"context"
	"fmt"
	"os"

	"app/cmd/dependency"
	"app/cmd/migrateCmd"
	"app/internal/module/healthcheck"
	"app/internal/module/users"
	"app/internal/shared/middleware"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	Version = "0.1"

	rootCmd = &cobra.Command{
		Use:     "server",
		Version: Version,
		Short:   "Запуск Http REST API",
		Run: func(cmd *cobra.Command, args []string) {
			app := fx.New(inject())

			err := app.Start(cmd.Context())
			if err != nil {
				return
			}

			defer func(app *fx.App, ctx context.Context) {
				err := app.Stop(ctx)
				if err != nil {
					fmt.Println(err)
				}
			}(app, cmd.Context())

			<-app.Done()
		},
	}
)

func inject() fx.Option {
	return fx.Options(
		//fx.NopLogger,
		fx.Provide(
			dependency.ProvideConfig,
		),
		fx.Invoke(dependency.ProvideLogger),

		fx.Provide(
			dependency.ProvideServer,
		),

		users.Module,

		//  healthcheck
		fx.Invoke(healthcheck.Controller),
		// загружаю мидлваре в приложение
		fx.Provide(middleware.NewMiddleware),
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
