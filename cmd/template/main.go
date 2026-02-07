package main

import (
	"flag"

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
		ProvideServerModule(),
	)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
