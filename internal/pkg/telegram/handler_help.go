package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// HelpHandler обрабатывает команду /help — список доступных команд.
type HelpHandler struct {
	log *slog.Logger
}

func NewHelpHandler(log *slog.Logger) *HelpHandler {
	return &HelpHandler{log: log}
}

func (h *HelpHandler) Options() []bot.Option {
	return []bot.Option{
		bot.WithMessageTextHandler("/help", bot.MatchTypeExact, h.handle),
	}
}

const helpText = `Доступные команды:

/start — приветствие и ваш Chat ID
/help — список команд
/remind — создать напоминание
/cancel — отменить текущее действие
/app — открыть Mini App`

func (h *HelpHandler) handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	}); err != nil {
		h.log.Error("failed to send /help response", slog.Any("err", err))
	}
}
