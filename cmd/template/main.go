package main

import (
	"context"
	"flag"
	"log"

	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/controller/template"
	appOtel "github.com/vovanwin/template/internal/pkg/otel"

	"go.uber.org/fx"
)

func inject(configDir string) fx.Option {
	// Загружаем конфиг до fx, чтобы использовать его при конструировании модулей
	cfg, err := config.Load(&config.LoadOptions{ConfigDir: configDir})
	if err != nil {
		log.Fatalf("загрузка конфига: %v", err)
	}

	return fx.Options(
		fx.Supply(cfg),
		fx.Provide(
			ProvideLogger,
			ProvideServerConfig,
			ProvideOtel,
			ProvidePgx,
		),

		// gRPC сервисы
		template.Module(),

		// Сервер (автоматически собирает все registrators)
		ProvideServerModule(cfg),

		// Graceful shutdown OTEL при остановке приложения
		fx.Invoke(func(lc fx.Lifecycle, provider *appOtel.Provider) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return provider.Shutdown(ctx)
				},
			})
		}),
	)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
