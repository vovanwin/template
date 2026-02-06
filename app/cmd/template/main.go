package main

import (
	"flag"

	"go.uber.org/fx"
)

func inject(configDir string) fx.Option {
	options := []fx.Option{
		fx.Provide(
			ProvideConfig(configDir),
			ProvideLogger,
			ProvidePgx,
		),
	}
	return fx.Options(options...)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
