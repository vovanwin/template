package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// TimePicker — inline-клавиатура для выбора времени (часы → минуты).
type TimePicker struct {
	prefix            string
	onSelect          TimePickerOnSelect
	callbackHandlerID string
}

// TimePickerOnSelect вызывается при выборе времени.
type TimePickerOnSelect func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, hour, minute int)

// NewTimePicker создаёт TimePicker и регистрирует callback-хэндлер на боте.
func NewTimePicker(b *bot.Bot, prefix string, onSelect TimePickerOnSelect) *TimePicker {
	tp := &TimePicker{
		prefix:   prefix,
		onSelect: onSelect,
	}
	tp.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, prefix, bot.MatchTypePrefix, tp.callback)
	return tp
}

// HoursKeyboard возвращает клавиатуру выбора часа (0–23, сетка 6×4).
func (tp *TimePicker) HoursKeyboard() *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton
	for r := 0; r < 6; r++ {
		var row []models.InlineKeyboardButton
		for c := 0; c < 4; c++ {
			h := r*4 + c
			row = append(row, models.InlineKeyboardButton{
				Text:         fmt.Sprintf("%02d", h),
				CallbackData: fmt.Sprintf("%s_h_%d", tp.prefix, h),
			})
		}
		rows = append(rows, row)
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// MinutesKeyboard возвращает клавиатуру выбора минут (шаг 5, сетка 3×4).
func (tp *TimePicker) MinutesKeyboard(hour int) *models.InlineKeyboardMarkup {
	var rows [][]models.InlineKeyboardButton
	for r := 0; r < 3; r++ {
		var row []models.InlineKeyboardButton
		for c := 0; c < 4; c++ {
			m := r*4 + c
			m *= 5 // шаг 5 минут: 0, 5, 10, ..., 55
			row = append(row, models.InlineKeyboardButton{
				Text:         fmt.Sprintf("%02d:%02d", hour, m),
				CallbackData: fmt.Sprintf("%s_m_%d_%d", tp.prefix, hour, m),
			})
		}
		rows = append(rows, row)
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func (tp *TimePicker) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	mes := update.CallbackQuery.Message

	// Формат: prefix_h_HOUR или prefix_m_HOUR_MINUTE
	trimmed := strings.TrimPrefix(data, tp.prefix+"_")
	parts := strings.Split(trimmed, "_")

	switch {
	case len(parts) == 2 && parts[0] == "h":
		// Выбран час — показываем минуты
		hour, err := strconv.Atoi(parts[1])
		if err != nil {
			return
		}
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      mes.Message.Chat.ID,
			MessageID:   mes.Message.ID,
			Text:        fmt.Sprintf("Час: %02d\nВыберите минуты:", hour),
			ReplyMarkup: tp.MinutesKeyboard(hour),
		})

	case len(parts) == 3 && parts[0] == "m":
		// Выбраны минуты — финал
		hour, err := strconv.Atoi(parts[1])
		if err != nil {
			return
		}
		minute, err := strconv.Atoi(parts[2])
		if err != nil {
			return
		}

		// Удаляем клавиатуру
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    mes.Message.Chat.ID,
			MessageID: mes.Message.ID,
			Text:      fmt.Sprintf("Время: %02d:%02d", hour, minute),
		})

		tp.onSelect(ctx, b, mes, hour, minute)
	}
}
