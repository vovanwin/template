package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"
	"github.com/go-telegram/ui/datepicker"
	"github.com/google/uuid"
	"github.com/vovanwin/template/internal/pkg/timezone"
	"github.com/vovanwin/template/internal/repository"
	"github.com/vovanwin/template/internal/service"
)

// FSM states для создания напоминания.
const (
	stateDefault          fsm.StateID = "default"
	stateWaitTitle        fsm.StateID = "waitTitle"
	stateWaitDescription  fsm.StateID = "waitDescription"
	stateWaitDate         fsm.StateID = "waitDate"
	stateWaitTime         fsm.StateID = "waitTime"
	stateWaitConfirmation fsm.StateID = "waitConfirmation"
)

// ReminderHandler обрабатывает команду /remind с FSM для многошагового диалога.
type ReminderHandler struct {
	reminderService *service.ReminderService
	userRepo        *repository.UserRepo
	fsm             *fsm.FSM
	log             *slog.Logger

	initOnce sync.Once
	dp       *datepicker.DatePicker
	tp       *TimePicker
}

func NewReminderHandler(
	reminderService *service.ReminderService,
	userRepo *repository.UserRepo,
	log *slog.Logger,
) *ReminderHandler {
	h := &ReminderHandler{
		reminderService: reminderService,
		userRepo:        userRepo,
		log:             log,
	}

	h.fsm = fsm.New(stateDefault, nil)

	return h
}

func (h *ReminderHandler) Options() []bot.Option {
	return []bot.Option{
		bot.WithMessageTextHandler("/remind", bot.MatchTypeExact, h.handleRemind),
		bot.WithMessageTextHandler("/cancel", bot.MatchTypeExact, h.handleCancel),
		bot.WithCallbackQueryDataHandler("ack_reminder:", bot.MatchTypePrefix, h.handleAckCallback),
		bot.WithCallbackQueryDataHandler("confirm_interval:", bot.MatchTypePrefix, h.handleConfirmIntervalCallback),
		bot.WithDefaultHandler(h.handleDefault),
	}
}

// initWidgets ленивая инициализация datepicker и timepicker.
// Хэндлеры создаются до *bot.Bot, поэтому виджеты инициализируются при первом вызове.
func (h *ReminderHandler) initWidgets(b *bot.Bot) {
	h.initOnce.Do(func() {
		h.dp = datepicker.New(b, h.onDateSelected,
			datepicker.Language("ru"),
			datepicker.WithPrefix("reminder_dp"),
			datepicker.From(time.Now().In(timezone.UserLocation)),
		)
		h.tp = NewTimePicker(b, "reminder_tp", h.onTimeSelected)
	})
}

// onDateSelected вызывается при выборе даты в datepicker.
func (h *ReminderHandler) onDateSelected(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, date time.Time) {
	chatID := mes.Message.Chat.ID
	userID := chatID // для личных чатов совпадают

	h.fsm.Set(userID, "date", date)
	h.fsm.Transition(userID, stateWaitTime)

	h.initWidgets(b)
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        fmt.Sprintf("Дата: %s\nВыберите час:", date.Format("02.01.2006")),
		ReplyMarkup: h.tp.HoursKeyboard(),
	}); err != nil {
		h.log.Error("failed to send timepicker", slog.Any("err", err))
	}
}

// onTimeSelected вызывается при выборе времени в timepicker.
func (h *ReminderHandler) onTimeSelected(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, hour, minute int) {
	chatID := mes.Message.Chat.ID
	userID := chatID

	dateVal, ok := h.fsm.Get(userID, "date")
	if !ok {
		h.sendError(ctx, b, chatID, "Ошибка: дата не найдена. Попробуйте /remind заново.")
		h.fsm.Reset(userID)
		return
	}
	date, _ := dateVal.(time.Time)

	// Собираем дату+время в таймзоне пользователя, затем конвертируем в UTC для хранения.
	localTime := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, timezone.UserLocation)
	remindAtUTC := localTime.UTC()

	h.fsm.Set(userID, "remind_at_utc", remindAtUTC)
	h.fsm.Set(userID, "local_time", localTime)
	h.fsm.Transition(userID, stateWaitConfirmation)

	// Показываем выбор подтверждения
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Нет", CallbackData: "confirm_interval:0"},
			},
			{
				{Text: "5 мин", CallbackData: "confirm_interval:5"},
				{Text: "10 мин", CallbackData: "confirm_interval:10"},
				{Text: "15 мин", CallbackData: "confirm_interval:15"},
			},
			{
				{Text: "30 мин", CallbackData: "confirm_interval:30"},
				{Text: "60 мин", CallbackData: "confirm_interval:60"},
			},
		},
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Требуется подтверждение получения? Выберите интервал повтора:",
		ReplyMarkup: keyboard,
	}); err != nil {
		h.log.Error("failed to send confirmation prompt", slog.Any("err", err))
	}
}

// handleRemind запускает FSM-диалог создания напоминания.
func (h *ReminderHandler) handleRemind(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	h.initWidgets(b)
	h.fsm.Transition(userID, stateWaitTitle)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Введите название напоминания:",
	}); err != nil {
		h.log.Error("failed to send remind prompt", slog.Any("err", err))
	}
}

// handleCancel сбрасывает FSM-состояние.
func (h *ReminderHandler) handleCancel(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	current := h.fsm.Current(userID)
	h.fsm.Reset(userID)

	if current != stateDefault {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Действие отменено.",
		}); err != nil {
			h.log.Error("failed to send cancel response", slog.Any("err", err))
		}
	}
}

// handleDefault маршрутизирует свободный ввод на основе текущего FSM-состояния.
func (h *ReminderHandler) handleDefault(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	state := h.fsm.Current(userID)

	switch state {
	case stateWaitTitle:
		h.fsm.Set(userID, "title", text)
		h.fsm.Transition(userID, stateWaitDescription)
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Введите описание (или отправьте «-» чтобы пропустить):",
		}); err != nil {
			h.log.Error("failed to send description prompt", slog.Any("err", err))
		}

	case stateWaitDescription:
		desc := text
		if desc == "-" {
			desc = ""
		}
		h.fsm.Set(userID, "description", desc)
		h.fsm.Transition(userID, stateWaitDate)

		h.initWidgets(b)
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatID,
			Text:        "Выберите дату напоминания:",
			ReplyMarkup: h.dp,
		}); err != nil {
			h.log.Error("failed to send datepicker", slog.Any("err", err))
		}

	case stateWaitDate:
		// Ожидаем нажатие на datepicker, текст игнорируем
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Пожалуйста, выберите дату из календаря выше.",
		}); err != nil {
			h.log.Error("failed to send date hint", slog.Any("err", err))
		}

	case stateWaitTime:
		// Ожидаем нажатие на timepicker, текст игнорируем
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Пожалуйста, выберите время из кнопок выше.",
		}); err != nil {
			h.log.Error("failed to send time hint", slog.Any("err", err))
		}

	case stateWaitConfirmation:
		// Ожидаем нажатие на inline-кнопку подтверждения
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Пожалуйста, выберите интервал из кнопок выше.",
		}); err != nil {
			h.log.Error("failed to send confirmation hint", slog.Any("err", err))
		}

	default:
		h.log.Info("unhandled message", slog.Int64("chat_id", chatID), slog.String("text", text))
	}
}

// handleConfirmIntervalCallback обрабатывает выбор интервала подтверждения.
func (h *ReminderHandler) handleConfirmIntervalCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	userID := chatID

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	if h.fsm.Current(userID) != stateWaitConfirmation {
		return
	}

	data := update.CallbackQuery.Data
	// Формат: confirm_interval:N
	var intervalMinutes int
	if _, err := fmt.Sscanf(data, "confirm_interval:%d", &intervalMinutes); err != nil {
		h.sendError(ctx, b, chatID, "Ошибка при обработке выбора.")
		h.fsm.Reset(userID)
		return
	}

	titleVal, _ := h.fsm.Get(userID, "title")
	descVal, _ := h.fsm.Get(userID, "description")
	remindAtVal, _ := h.fsm.Get(userID, "remind_at_utc")
	localTimeVal, _ := h.fsm.Get(userID, "local_time")
	h.fsm.Reset(userID)

	title, _ := titleVal.(string)
	desc, _ := descVal.(string)
	remindAtUTC, _ := remindAtVal.(time.Time)
	localTime, _ := localTimeVal.(time.Time)

	requireConfirmation := intervalMinutes > 0

	h.finishReminder(ctx, b, chatID, title, desc, remindAtUTC, localTime, requireConfirmation, intervalMinutes)
}

// handleAckCallback обрабатывает нажатие кнопки "Подтвердить" на уведомлении.
func (h *ReminderHandler) handleAckCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	data := update.CallbackQuery.Data

	// Формат: ack_reminder:<reminder_id>
	var reminderIDStr string
	if len(data) > len("ack_reminder:") {
		reminderIDStr = data[len("ack_reminder:"):]
	}

	if reminderIDStr == "" {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Ошибка: неверные данные",
		})
		return
	}

	// Ищем пользователя по chat_id
	user, err := h.userRepo.GetByChatID(ctx, chatID)
	if err != nil || user == nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Ошибка: пользователь не найден",
		})
		return
	}

	reminderID, err := uuid.Parse(reminderIDStr)
	if err != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Ошибка: неверный ID",
		})
		return
	}

	if err := h.reminderService.AcknowledgeReminder(ctx, user.ID, reminderID); err != nil {
		h.log.Error("failed to acknowledge reminder", slog.Any("err", err))
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Ошибка подтверждения",
		})
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "✅ Подтверждено!",
	})
}

// finishReminder завершает создание напоминания.
func (h *ReminderHandler) finishReminder(ctx context.Context, b *bot.Bot, chatID int64, title, desc string, remindAtUTC, localTime time.Time, requireConfirmation bool, repeatIntervalMinutes int) {
	user, err := h.userRepo.GetByChatID(ctx, chatID)
	if err != nil {
		h.log.Error("failed to find user by chat_id", slog.Any("err", err), slog.Int64("chat_id", chatID))
		h.sendError(ctx, b, chatID, "Ошибка при поиске пользователя.")
		return
	}
	if user == nil {
		h.sendError(ctx, b, chatID, "Ваш аккаунт не привязан. Укажите Chat ID в настройках профиля.")
		return
	}

	rem, err := h.reminderService.CreateReminder(ctx, user.ID, title, desc, remindAtUTC, chatID, requireConfirmation, repeatIntervalMinutes)
	if err != nil {
		h.log.Error("failed to create reminder", slog.Any("err", err))
		h.sendError(ctx, b, chatID, "Ошибка при создании напоминания.")
		return
	}

	// Показываем пользователю время в его таймзоне
	displayTime := timezone.FormatUser(rem.RemindAt, "02.01.2006 15:04")
	msg := fmt.Sprintf("Напоминание создано!\n\nНазвание: %s\nВремя: %s", rem.Title, displayTime)
	if desc != "" {
		msg = fmt.Sprintf("Напоминание создано!\n\nНазвание: %s\nОписание: %s\nВремя: %s", rem.Title, desc, displayTime)
	}
	if requireConfirmation {
		msg += fmt.Sprintf("\nПодтверждение: каждые %d мин", repeatIntervalMinutes)
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	}); err != nil {
		h.log.Error("failed to send reminder created message", slog.Any("err", err))
	}
}

func (h *ReminderHandler) sendError(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		h.log.Error("failed to send error message", slog.Any("err", err))
	}
}
