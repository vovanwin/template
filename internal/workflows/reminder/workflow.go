package reminder

import (
	"time"

	"github.com/vovanwin/template/internal/model"
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
		req:         input.Req,
		cancel:      input.CancelReminder,
		acknowledge: input.AcknowledgeReminder,
		status:      model.ReminderStatusPending,
	}, nil
}

type scheduleReminderWorkflow struct {
	req         *reminderv1.ScheduleReminderRequest
	cancel      *reminderv1.CancelReminderSignal
	acknowledge *reminderv1.AcknowledgeReminderSignal
	status      model.ReminderStatus
}

func (w *scheduleReminderWorkflow) Execute(ctx workflow.Context) (*reminderv1.ScheduleReminderResponse, error) {
	log := workflow.GetLogger(ctx)
	reminderID := w.req.GetReminderId()
	workflowID := workflow.GetInfo(ctx).WorkflowExecution.ID

	// Retry-политики и таймауты заданы в proto (reminder.proto) для каждой activity:
	//   SendTelegramNotification: start_to_close=30s, max_attempts=5
	//   UpdateReminderStatus:     start_to_close=10s, max_attempts=10
	// Proto-сгенерированные хелперы автоматически применяют эти настройки.

	log.Info("reminder workflow started",
		"reminder_id", reminderID,
		"workflow_id", workflowID,
		"remind_at", w.req.GetRemindAt().AsTime(),
	)

	// 1. Обновляем статус в БД на "processing"
	w.setStatus(ctx, reminderID, model.ReminderStatusProcessing)

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
	var timerErr error

	sel.AddFuture(timerFuture, func(f workflow.Future) {
		timerErr = f.Get(timerCtx, nil)
	})

	sel.AddReceive(w.cancel.Channel, func(ch workflow.ReceiveChannel, more bool) {
		cancelled = true
		timerCancel()
	})

	sel.Select(ctx)

	// 2. Обработка отмены
	if cancelled {
		log.Info("reminder cancelled", "reminder_id", reminderID)
		w.setStatus(ctx, reminderID, model.ReminderStatusCancelled)
		return w.response(workflowID), nil
	}

	// Таймер отменён не через наш сигнал (например terminate из Temporal UI)
	if timerErr != nil {
		log.Error("timer error", "error", timerErr, "reminder_id", reminderID)
		w.setStatus(ctx, reminderID, model.ReminderStatusFailed)
		return w.response(workflowID), nil
	}

	// 3. Проверяем chatID
	chatID := w.req.GetTelegramChatId()
	if chatID == 0 {
		log.Error("telegram chat_id is 0, cannot send notification", "reminder_id", reminderID)
		w.setStatus(ctx, reminderID, model.ReminderStatusFailed)
		return w.response(workflowID), nil
	}

	// 4. Отправляем уведомление в Telegram
	// При ошибке proto-хелпер сделает до 5 retry (заданы в proto).
	// Если все попытки провалились — workflow завершается со статусом "failed",
	// а НЕ возвращает error (иначе Temporal бесконечно ретраит весь workflow).
	err := reminderv1.SendTelegramNotification(ctx, &reminderv1.SendTelegramNotificationRequest{
		ChatId:              chatID,
		Title:               w.req.GetTitle(),
		Description:         w.req.GetDescription(),
		RequireConfirmation: w.req.GetRequireConfirmation(),
		ReminderId:          reminderID,
	})
	if err != nil {
		log.Error("failed to send telegram notification after retries",
			"error", err,
			"reminder_id", reminderID,
		)
		w.setStatus(ctx, reminderID, model.ReminderStatusFailed)
		return w.response(workflowID), nil
	}

	// 5. Если требуется подтверждение — цикл повторных уведомлений
	if w.req.GetRequireConfirmation() && w.req.GetRepeatIntervalMinutes() > 0 {
		repeatInterval := time.Duration(w.req.GetRepeatIntervalMinutes()) * time.Minute
		maxRetryDuration := 10 * time.Hour
		startTime := workflow.Now(ctx)

		for {
			elapsed := workflow.Now(ctx).Sub(startTime)
			if elapsed >= maxRetryDuration {
				log.Info("confirmation timeout exceeded, marking as sent", "reminder_id", reminderID)
				break
			}

			// Ждём repeat_interval, AcknowledgeReminder или CancelReminder
			repeatTimerCtx, repeatTimerCancel := workflow.WithCancel(ctx)
			repeatFuture := workflow.NewTimer(repeatTimerCtx, repeatInterval)

			repeatSel := workflow.NewSelector(ctx)

			var acknowledged, repeatCancelled bool
			var repeatTimerErr error

			repeatSel.AddFuture(repeatFuture, func(f workflow.Future) {
				repeatTimerErr = f.Get(repeatTimerCtx, nil)
			})

			repeatSel.AddReceive(w.acknowledge.Channel, func(ch workflow.ReceiveChannel, more bool) {
				acknowledged = true
				repeatTimerCancel()
			})

			repeatSel.AddReceive(w.cancel.Channel, func(ch workflow.ReceiveChannel, more bool) {
				repeatCancelled = true
				repeatTimerCancel()
			})

			repeatSel.Select(ctx)

			if acknowledged {
				log.Info("reminder acknowledged", "reminder_id", reminderID)
				w.setStatus(ctx, reminderID, model.ReminderStatusSent)
				return w.response(workflowID), nil
			}

			if repeatCancelled {
				log.Info("reminder cancelled during confirmation loop", "reminder_id", reminderID)
				w.setStatus(ctx, reminderID, model.ReminderStatusCancelled)
				return w.response(workflowID), nil
			}

			if repeatTimerErr != nil {
				log.Error("repeat timer error", "error", repeatTimerErr, "reminder_id", reminderID)
				break
			}

			// Повторная отправка уведомления
			err := reminderv1.SendTelegramNotification(ctx, &reminderv1.SendTelegramNotificationRequest{
				ChatId:              chatID,
				Title:               w.req.GetTitle(),
				Description:         w.req.GetDescription(),
				RequireConfirmation: true,
				ReminderId:          reminderID,
			})
			if err != nil {
				log.Error("failed to resend notification", "error", err, "reminder_id", reminderID)
			}
		}
	}

	// 6. Обновляем статус в БД на "sent"
	log.Info("reminder sent", "reminder_id", reminderID, "chat_id", chatID)
	w.setStatus(ctx, reminderID, model.ReminderStatusSent)

	return w.response(workflowID), nil
}

// setStatus обновляет статус в workflow и в БД через activity.
func (w *scheduleReminderWorkflow) setStatus(ctx workflow.Context, reminderID string, status model.ReminderStatus) {
	w.status = status
	if reminderID == "" {
		return
	}
	err := reminderv1.UpdateReminderStatus(ctx, &reminderv1.UpdateReminderStatusRequest{
		ReminderId: reminderID,
		Status:     status.String(),
	})
	if err != nil {
		workflow.GetLogger(ctx).Error("failed to update reminder status",
			"error", err,
			"reminder_id", reminderID,
			"status", status.String(),
		)
	}
}

// response формирует ответ с текущим статусом.
func (w *scheduleReminderWorkflow) response(workflowID string) *reminderv1.ScheduleReminderResponse {
	return &reminderv1.ScheduleReminderResponse{
		WorkflowId: workflowID,
		Status:     w.status.String(),
	}
}

func (w *scheduleReminderWorkflow) GetReminderStatus() (*reminderv1.GetReminderStatusResponse, error) {
	return &reminderv1.GetReminderStatusResponse{
		Status: w.status.String(),
	}, nil
}

// ReminderDuration вычисляет задержку до напоминания.
func ReminderDuration(remindAt time.Time) time.Duration {
	d := time.Until(remindAt)
	if d < 0 {
		return 0
	}
	return d
}
