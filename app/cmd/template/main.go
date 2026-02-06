package main

import (
	"flag"
	"log/slog"

	"go.uber.org/fx"
)

func inject(configDir string) fx.Option {
	options := []fx.Option{
		fx.Provide(
			ProvideConfig(configDir),
			ProvideLogger,
			ProvidePgx,
		), fx.Invoke(func(log *slog.Logger) {
			log.Info("test")
		}),
	}
	return fx.Options(options...)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
