package telegram

import "github.com/go-telegram/bot"

// HandlerRegistrar — интерфейс для регистрации обработчиков команд Telegram бота.
// Каждый хэндлер возвращает опции, которые передаются в bot.New().
type HandlerRegistrar interface {
	Options() []bot.Option
}
