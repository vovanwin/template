package reminder

import (
	"fmt"
	"time"

	reminderv1 "github.com/vovanwin/template/pkg/temporal/reminder"
	"go.temporal.io/sdk/workflow"
)

// Workflows реализует интерфейс ReminderWorkflows.
type Workflows struct{}

func NewWorkflows() *Workflows {
	return &Workflows{}
}

// ScheduleReminder запускает workflow: ждёт до remind_at, затем отправляет уведомление.
func (w *Workflows) ScheduleReminder(ctx workflow.Context, input *reminderv1.ScheduleReminderWorkflowInput) (reminderv1.ScheduleReminderWorkflow, error) {
	return &scheduleReminderWorkflow{
		req:    input.Req,
		cancel: input.CancelReminder,
		status: "pending",
	}, nil
}

type scheduleReminderWorkflow struct {
	req    *reminderv1.ScheduleReminderRequest
	cancel *reminderv1.CancelReminderSignal
	status string
}

func (w *scheduleReminderWorkflow) Execute(ctx workflow.Context) (*reminderv1.ScheduleReminderResponse, error) {
	remindAt := w.req.GetRemindAt().AsTime()
	now := workflow.Now(ctx)
	duration := remindAt.Sub(now)
	if duration < 0 {
		duration = 0
	}

	// Создаём таймер и канал отмены
	timerCtx, timerCancel := workflow.WithCancel(ctx)
	timerFuture := workflow.NewTimer(timerCtx, duration)

	// Selector для отслеживания таймера и сигнала отмены
	sel := workflow.NewSelector(ctx)

	var cancelled bool

	sel.AddFuture(timerFuture, func(f workflow.Future) {
		// Таймер сработал — отправляем уведомление
	})

	sel.AddReceive(w.cancel.Channel, func(ch workflow.ReceiveChannel, more bool) {
		// Сигнал отмены
		cancelled = true
		timerCancel()
	})

	sel.Select(ctx)

	if cancelled {
		w.status = "cancelled"
		return &reminderv1.ScheduleReminderResponse{
			WorkflowId: workflow.GetInfo(ctx).WorkflowExecution.ID,
			Status:     "cancelled",
		}, nil
	}

	// Отправляем уведомление в Telegram
	chatID := w.req.GetTelegramChatId()
	if chatID != 0 {
		title := w.req.GetTitle()
		description := w.req.GetDescription()
		err := reminderv1.SendTelegramNotification(ctx, &reminderv1.SendTelegramNotificationRequest{
			ChatId:      chatID,
			Title:       title,
			Description: description,
		})
		if err != nil {
			w.status = "failed"
			return nil, fmt.Errorf("send telegram notification: %w", err)
		}
	}

	w.status = "sent"
	return &reminderv1.ScheduleReminderResponse{
		WorkflowId: workflow.GetInfo(ctx).WorkflowExecution.ID,
		Status:     "sent",
	}, nil
}

func (w *scheduleReminderWorkflow) GetReminderStatus() (*reminderv1.GetReminderStatusResponse, error) {
	return &reminderv1.GetReminderStatusResponse{
		Status: w.status,
	}, nil
}

// WorkflowName для внешнего использования.
func WorkflowName() string {
	return reminderv1.ScheduleReminderWorkflowName
}

// TaskQueue для внешнего использования.
func TaskQueue() string {
	return "reminder-v1"
}

// ReminderDuration вычисляет задержку до напоминания.
func ReminderDuration(remindAt time.Time) time.Duration {
	d := time.Until(remindAt)
	if d < 0 {
		return 0
	}
	return d
}
