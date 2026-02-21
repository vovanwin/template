package telegram

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/config"
	"github.com/vovanwin/template/internal/repository"
	"github.com/vovanwin/template/internal/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// HandlerParams — параметры для сбора всех хэндлеров через fx.
type HandlerParams struct {
	fx.In
	Handlers []HandlerRegistrar `group:"telegram_handlers"`
}

// Module возвращает fx.Option для Telegram бота с модульными хэндлерами.
func Module() fx.Option {
	return fx.Module("telegram",
		// Хэндлеры
		fx.Provide(
			fx.Annotate(
				func(log *slog.Logger) HandlerRegistrar {
					return NewStartHandler(log)
				},
				fx.ResultTags(`group:"telegram_handlers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(log *slog.Logger) HandlerRegistrar {
					return NewHelpHandler(log)
				},
				fx.ResultTags(`group:"telegram_handlers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(cfg *config.Config, log *slog.Logger) HandlerRegistrar {
					return NewMiniAppHandler(cfg.Telegram.MiniAppURL, log)
				},
				fx.ResultTags(`group:"telegram_handlers"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				func(
					reminderService *service.ReminderService,
					userRepo *repository.UserRepo,
					log *slog.Logger,
				) HandlerRegistrar {
					return NewReminderHandler(reminderService, userRepo, log)
				},
				fx.ResultTags(`group:"telegram_handlers"`),
			),
		),

		// Bot
		fx.Provide(func(cfg *config.Config, log *slog.Logger, hp HandlerParams) (*Bot, error) {
			var opts []bot.Option
			for _, h := range hp.Handlers {
				opts = append(opts, h.Options()...)
			}
			return New(cfg.Telegram.Token, cfg.Telegram.WebhookURL, log, opts...)
		}),

		// Lifecycle: start/stop бота
		fx.Invoke(func(lc fx.Lifecycle, tgBot *Bot) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go tgBot.Start(ctx)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					tgBot.Stop()
					return nil
				},
			})
		}),

		// GatewayRegistrator для webhook-роута
		fx.Provide(
			fx.Annotate(
				func(tgBot *Bot) server.GatewayRegistrator {
					return func(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
						if !tgBot.IsWebhookMode() {
							return nil
						}
						return mux.HandlePath("POST", "/api/telegram/webhook", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
							tgBot.WebhookHandler()(w, r)
						})
					}
				},
				fx.ResultTags(`group:"gateway_registrators"`),
			),
		),
	)
}
