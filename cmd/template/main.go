package main

import (
	"flag"

	"github.com/vovanwin/template/internal/controller/template"

	"go.uber.org/fx"
)

func inject(configDir string) fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideConfig(configDir),
			ProvideLogger,
			ProvideServerConfig,
			ProvidePgx,
		),

		// gRPC сервисы
		template.Module(),

		// Сервер (автоматически собирает все registrators)
		ProvideServerModule(),
	)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
