package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// MiniAppHandler обрабатывает команду /app — отправляет кнопку для открытия WebApp.
type MiniAppHandler struct {
	miniAppURL string
	log        *slog.Logger
}

func NewMiniAppHandler(miniAppURL string, log *slog.Logger) *MiniAppHandler {
	return &MiniAppHandler{miniAppURL: miniAppURL, log: log}
}

func (h *MiniAppHandler) Options() []bot.Option {
	return []bot.Option{
		bot.WithMessageTextHandler("/app", bot.MatchTypeExact, h.handle),
	}
}

func (h *MiniAppHandler) handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID

	if h.miniAppURL == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Mini App не настроен.",
		}); err != nil {
			h.log.Error("failed to send /app response", slog.Any("err", err))
		}
		return
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Откройте приложение:",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text:   "Открыть Mini App",
						WebApp: &models.WebAppInfo{URL: h.miniAppURL},
					},
				},
			},
		},
	}); err != nil {
		h.log.Error("failed to send /app response", slog.Any("err", err))
	}
}
