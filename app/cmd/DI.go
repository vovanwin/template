package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/vovanwin/template/app/cmd/dependency"
	"github.com/vovanwin/template/app/cmd/migrateCmd"
	"github.com/vovanwin/template/app/internal/module/users"
	"github.com/vovanwin/template/app/internal/module/web"
	"github.com/vovanwin/template/app/internal/shared/middleware"

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
			dependency.ProvidePgx,
			dependency.ProvidePool,
			dependency.ProvideSessionManager,
		),
		fx.Invoke(dependency.ProvideLogger),

		fx.Provide(
			dependency.ProvideServer,
		),

		// start additional servers via lifecycle hooks
		fx.Invoke(
			dependency.ProvideDebugServer,
			dependency.ProvideSwaggerServer,
			dependency.ProvideGRPCServer,
		),

		users.Module,
		web.Module,

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
