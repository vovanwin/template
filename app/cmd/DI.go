package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/vovanwin/template/app/cmd/dependency"
	"github.com/vovanwin/template/app/cmd/migrateCmd"
	"github.com/vovanwin/template/app/internal/module/users"
	"github.com/vovanwin/template/app/internal/shared/middleware"
	"github.com/vovanwin/template/app/internal/workflows"

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
	// Загружаем конфигурацию сначала, чтобы знать какие компоненты включены
	config, err := dependency.ProvideConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	options := []fx.Option{
		//fx.NopLogger,
		fx.Supply(config), // Предоставляем конфигурацию как уже созданную зависимость
		fx.Provide(
			dependency.InitLogger,
			dependency.ProvidePgx,
			dependency.ProvidePool,
			dependency.ProvideJWTService,
			dependency.ProvideTemporal, // Всегда предоставляем Temporal (может быть nil)
		),
		// Логируем статус компонентов при старте
		fx.Invoke(func() { dependency.LogComponentsStatus(config) }),
	}

	// Всегда добавляем провайдер HTTP сервера (он может быть отключен внутри)
	options = append(options, fx.Provide(dependency.ProvideServer))

	// Условно добавляем дополнительные серверы
	var invokeOptions []interface{}

	if config.Server.EnableDebug {
		invokeOptions = append(invokeOptions, dependency.ProvideDebugServer)
	}

	if config.Server.EnableSwagger {
		invokeOptions = append(invokeOptions, dependency.ProvideSwaggerServer)
	}

	if config.Server.EnableGRPC {
		invokeOptions = append(invokeOptions, dependency.ProvideGRPCServer)
	}

	if len(invokeOptions) > 0 {
		options = append(options, fx.Invoke(invokeOptions...))
	}

	// Добавляем модули приложения
	options = append(options,
		users.Module,
		fx.Provide(middleware.NewMiddleware),
	)

	// Условно добавляем workflows модуль только если Temporal включен
	if config.Server.EnableTemporal {
		options = append(options, workflows.Module)
	}

	return fx.Options(options...)
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
