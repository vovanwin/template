package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// StartHandler обрабатывает команду /start — приветствие с Chat ID.
type StartHandler struct {
	log *slog.Logger
}

func NewStartHandler(log *slog.Logger) *StartHandler {
	return &StartHandler{log: log}
}

func (h *StartHandler) Options() []bot.Option {
	return []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, h.handle),
	}
}

func (h *StartHandler) handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	msg := fmt.Sprintf("Привет! Ваш Chat ID: %d\nУкажите его в настройках профиля для получения напоминаний.", chatID)
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	}); err != nil {
		h.log.Error("failed to send /start response", slog.Any("err", err))
	}
}
