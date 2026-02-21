package reminder

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vovanwin/template/internal/pkg/telegram"
	"github.com/vovanwin/template/internal/repository"
	reminderv1 "github.com/vovanwin/template/pkg/temporal/reminder"
)

// Activities реализует интерфейс ReminderActivities.
type Activities struct {
	bot  *telegram.Bot
	repo *repository.ReminderRepo
}

func NewActivities(bot *telegram.Bot, repo *repository.ReminderRepo) *Activities {
	return &Activities{bot: bot, repo: repo}
}

// SendTelegramNotification отправляет уведомление о напоминании в Telegram.
func (a *Activities) SendTelegramNotification(ctx context.Context, req *reminderv1.SendTelegramNotificationRequest) error {
	text := fmt.Sprintf("🔔 %s", req.GetTitle())
	if desc := req.GetDescription(); desc != "" {
		text += fmt.Sprintf("\n\n%s", desc)
	}
	return a.bot.SendMessage(ctx, req.GetChatId(), text)
}

// UpdateReminderStatus обновляет статус напоминания в базе данных.
func (a *Activities) UpdateReminderStatus(ctx context.Context, req *reminderv1.UpdateReminderStatusRequest) error {
	id, err := uuid.Parse(req.GetReminderId())
	if err != nil {
		return fmt.Errorf("parse reminder id: %w", err)
	}
	return a.repo.UpdateStatus(ctx, id, req.GetStatus())
}
