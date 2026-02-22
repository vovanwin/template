package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/vovanwin/template/internal/model"
	"github.com/vovanwin/template/internal/pkg/temporal"
	"github.com/vovanwin/template/internal/repository"
	reminderv1 "github.com/vovanwin/template/pkg/temporal/reminder"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReminderService struct {
	repo     *repository.ReminderRepo
	temporal *temporal.Service
	log      *slog.Logger
}

func NewReminderService(
	repo *repository.ReminderRepo,
	temporalSvc *temporal.Service,
	log *slog.Logger,
) *ReminderService {
	return &ReminderService{
		repo:     repo,
		temporal: temporalSvc,
		log:      log,
	}
}

func (s *ReminderService) CreateReminder(ctx context.Context, userID uuid.UUID, title, description string, remindAt time.Time, telegramChatID int64, requireConfirmation bool, repeatIntervalMinutes int) (*repository.Reminder, error) {
	rem, err := s.repo.Create(ctx, userID, title, description, remindAt, requireConfirmation, repeatIntervalMinutes)
	if err != nil {
		return nil, fmt.Errorf("create reminder in db: %w", err)
	}

	// Запускаем Temporal workflow на очереди из proto (reminder-v1)
	workflowID := fmt.Sprintf("reminder/%s", rem.ID.String())
	opts := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: reminderv1.ReminderTaskQueue,
	}

	req := &reminderv1.ScheduleReminderRequest{
		ReminderId:            rem.ID.String(),
		UserId:                userID.String(),
		Title:                 title,
		Description:           description,
		RemindAt:              timestamppb.New(remindAt),
		TelegramChatId:        telegramChatID,
		RequireConfirmation:   requireConfirmation,
		RepeatIntervalMinutes: int32(repeatIntervalMinutes),
	}

	run, err := s.temporal.GetClient().ExecuteWorkflow(ctx, opts, reminderv1.ScheduleReminderWorkflowName, req)
	if err != nil {
		s.log.Error("failed to start reminder workflow", slog.Any("err", err), slog.String("reminder_id", rem.ID.String()))
		// Не возвращаем ошибку — напоминание создано в БД, workflow можно перезапустить
	} else {
		_ = s.repo.UpdateWorkflowID(ctx, rem.ID, run.GetID())
		rem.WorkflowID = run.GetID()
	}

	return rem, nil
}

func (s *ReminderService) ListReminders(ctx context.Context, userID uuid.UUID) ([]repository.Reminder, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *ReminderService) ListRemindersPaged(ctx context.Context, userID uuid.UUID, page, pageSize int) (*repository.PagedReminders, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListByUserIDPaged(ctx, userID, page, pageSize)
}

func (s *ReminderService) CancelReminder(ctx context.Context, userID, reminderID uuid.UUID) error {
	rem, err := s.repo.GetByID(ctx, reminderID)
	if err != nil {
		return fmt.Errorf("get reminder: %w", err)
	}
	if rem == nil {
		return fmt.Errorf("reminder not found")
	}
	if rem.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	// Отправляем сигнал отмены в workflow
	if rem.WorkflowID != "" {
		err = s.temporal.GetClient().GetClient().SignalWorkflow(ctx, rem.WorkflowID, "", reminderv1.CancelReminderSignalName, nil)
		if err != nil {
			s.log.Warn("failed to cancel workflow", slog.Any("err", err), slog.String("workflow_id", rem.WorkflowID))
		}
	}

	return s.repo.UpdateStatus(ctx, reminderID, model.ReminderStatusCancelled.String())
}

func (s *ReminderService) AcknowledgeReminder(ctx context.Context, userID, reminderID uuid.UUID) error {
	rem, err := s.repo.GetByID(ctx, reminderID)
	if err != nil {
		return fmt.Errorf("get reminder: %w", err)
	}
	if rem == nil {
		return fmt.Errorf("reminder not found")
	}
	if rem.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	if rem.WorkflowID != "" {
		err = s.temporal.GetClient().GetClient().SignalWorkflow(ctx, rem.WorkflowID, "", reminderv1.AcknowledgeReminderSignalName, nil)
		if err != nil {
			s.log.Warn("failed to acknowledge workflow", slog.Any("err", err), slog.String("workflow_id", rem.WorkflowID))
			return fmt.Errorf("acknowledge workflow: %w", err)
		}
	}

	return nil
}

func (s *ReminderService) DeleteReminder(ctx context.Context, userID, reminderID uuid.UUID) error {
	rem, err := s.repo.GetByID(ctx, reminderID)
	if err != nil {
		return fmt.Errorf("get reminder: %w", err)
	}
	if rem == nil {
		return fmt.Errorf("reminder not found")
	}
	if rem.UserID != userID {
		return fmt.Errorf("forbidden")
	}

	// Если workflow активен — пробуем отменить
	if rem.WorkflowID != "" && rem.Status == model.ReminderStatusPending.String() {
		_ = s.temporal.GetClient().GetClient().SignalWorkflow(ctx, rem.WorkflowID, "", reminderv1.CancelReminderSignalName, nil)
	}

	return s.repo.Delete(ctx, reminderID)
}
