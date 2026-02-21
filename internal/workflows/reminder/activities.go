package reminder

import (
	"context"
	"fmt"

	"github.com/vovanwin/template/internal/pkg/telegram"
	reminderv1 "github.com/vovanwin/template/pkg/temporal/reminder"
)

// Activities реализует интерфейс ReminderActivities.
type Activities struct {
	bot *telegram.Bot
}

func NewActivities(bot *telegram.Bot) *Activities {
	return &Activities{bot: bot}
}

// SendTelegramNotification отправляет уведомление о напоминании в Telegram.
func (a *Activities) SendTelegramNotification(ctx context.Context, req *reminderv1.SendTelegramNotificationRequest) error {
	text := fmt.Sprintf("🔔 %s", req.GetTitle())
	if desc := req.GetDescription(); desc != "" {
		text += fmt.Sprintf("\n\n%s", desc)
	}
	return a.bot.SendMessage(ctx, req.GetChatId(), text)
}
