package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"time"

	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/controller/auth"
	"github.com/vovanwin/template/internal/controller/template"
	"github.com/vovanwin/template/internal/controller/ui"
	"github.com/vovanwin/template/internal/pkg/events"
	"github.com/vovanwin/template/internal/pkg/jwt"
	"github.com/vovanwin/template/internal/pkg/telegram"
	"github.com/vovanwin/template/internal/pkg/temporal"
	"github.com/vovanwin/template/internal/repository"
	"github.com/vovanwin/template/internal/service"
	"github.com/vovanwin/template/internal/workflows"

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
			func(cfg *config.Config, log *slog.Logger) (*temporal.Service, error) {
				return temporal.NewService(temporal.ServiceConfig{
					Client: temporal.Config{
						Host:      cfg.Temporal.Host,
						Port:      cfg.Temporal.Port,
						Namespace: cfg.Temporal.Namespace,
					},
					Worker: temporal.WorkerConfig{TaskQueue: cfg.Temporal.TaskQueue},
				}, log.With("component", "temporal"))
			},
			ProvideServerConfig,
			ProvidePgx,
			repository.NewUserRepo,
			repository.NewSessionRepo,
			repository.NewReminderRepo,
			service.NewAuthService,
			service.NewReminderService,
			func() jwt.JWTService {
				return jwtService
			},
		),
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.StopHook(closeFn))
		}),
		fx.Invoke(func(lc fx.Lifecycle, temporalSvc *temporal.Service) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return temporalSvc.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					temporalSvc.Stop(ctx)
					return nil
				},
			})
		}),

		// Telegram бот (модульная архитектура)
		telegram.Module(),

		// gRPC сервисы
		template.Module(),
		auth.Module(),
		ui.Module(),

		// Workflows
		workflows.Module,

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

func init() {
	time.Local = time.UTC
}
