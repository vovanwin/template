package main

import (
	"flag"
	"log"

	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/controller/auth"
	"github.com/vovanwin/template/internal/controller/template"
	"github.com/vovanwin/template/internal/controller/ui"
	"github.com/vovanwin/template/internal/pkg/events"
	"github.com/vovanwin/template/internal/pkg/jwt"
	"github.com/vovanwin/template/internal/repository"
	"github.com/vovanwin/template/internal/service"

	"go.uber.org/fx"
)

func inject(configDir string) fx.Option {
	// Загружаем конфиг до fx, чтобы использовать его при конструировании модулей
	cfg, err := config.Load(&config.LoadOptions{ConfigDir: configDir})
	if err != nil {
		log.Fatalf("загрузка конфига: %v", err)
	}

	flags, closeFn := ProvideFlags(cfg)
	jwtService := jwt.NewJWTService(cfg.JWT.SignKey, cfg.JWT.TokenTtl, cfg.JWT.RefreshTokenTtl)
	eventBus := events.NewBus()

	return fx.Options(
		fx.Supply(cfg),
		fx.Supply(flags),
		fx.Supply(eventBus),
		fx.Provide(
			ProvideLogger,
			ProvideServerConfig,
			ProvidePgx,
			repository.NewUserRepo,
			repository.NewSessionRepo,
			service.NewAuthService,
			func() jwt.JWTService {
				return jwtService
			},
		),
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.StopHook(closeFn))
		}),

		// gRPC сервисы
		template.Module(),
		auth.Module(),
		ui.Module(),

		// Сервер (автоматически собирает все registrators)
		ProvideServerModule(cfg, flags, jwtService),
	)
}

func main() {
	configDir := flag.String("config", "./app/config", "путь к директории с конфигами")
	flag.Parse()

	app := fx.New(inject(*configDir))

	app.Run()
}
