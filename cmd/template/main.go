package main

import (
	"flag"
	"log"

	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/controller/template"

	"go.uber.org/fx"
)

func inject(configDir string) fx.Option {
	// Загружаем конфиг до fx, чтобы использовать его при конструировании модулей
	cfg, err := config.Load(&config.LoadOptions{ConfigDir: configDir})
	if err != nil {
		log.Fatalf("загрузка конфига: %v", err)
	}

	flags := ProvideFlags()

	return fx.Options(
		fx.Supply(cfg),
		fx.Supply(flags),
		fx.Provide(
			ProvideLogger,
			ProvideServerConfig,
			ProvidePgx,
		),

		// gRPC сервисы
		template.Module(),

		// Сервер (автоматически собирает все registrators)
		ProvideServerModule(cfg, flags),
	)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
