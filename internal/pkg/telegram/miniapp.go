package telegram

import (
	"net/url"

	"github.com/go-telegram/bot"
)

// ValidateWebAppRequest валидирует запрос из Telegram Mini App (WebApp).
// Проверяет подпись initData, которую Telegram передаёт при открытии WebApp.
func ValidateWebAppRequest(token string, values url.Values) (*bot.WebAppUser, bool) {
	return bot.ValidateWebappRequest(values, token)
}
