package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Bot обёртка над Telegram ботом для отправки уведомлений.
// Поддерживает два режима: long polling (локальная разработка) и webhook (продакшен).
type Bot struct {
	bot        *bot.Bot
	token      string
	webhookURL string
	cancel     context.CancelFunc
	log        *slog.Logger
}

// New создаёт нового Telegram бота. Если token пустой — возвращает noop-бота.
// webhookURL определяет режим работы:
//   - пустой — long polling (для локальной разработки)
//   - непустой — webhook (для продакшена)
func New(token, webhookURL string, log *slog.Logger, opts ...bot.Option) (*Bot, error) {
	if token == "" {
		log.Warn("telegram bot token is empty, notifications disabled")
		return &Bot{log: log}, nil
	}

	// HTTP timeout должен быть больше long polling timeout,
	// иначе HTTP-клиент оборвёт запрос раньше, чем Telegram ответит.
	httpClient := &http.Client{Timeout: 60 * time.Second}
	defaultOpts := []bot.Option{
		bot.WithHTTPClient(30*time.Second, httpClient),
	}
	allOpts := append(defaultOpts, opts...)

	b, err := bot.New(token, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	return &Bot{bot: b, token: token, webhookURL: webhookURL, log: log}, nil
}

// Token возвращает токен бота (нужен для ValidateWebappRequest).
func (b *Bot) Token() string {
	return b.token
}

// Start запускает бота в нужном режиме:
//   - long polling если webhookURL пустой
//   - webhook если webhookURL задан (нужно дополнительно подключить WebhookHandler к роутеру)
//
// Используем context.Background() вместо переданного контекста,
// потому что fx OnStart контекст имеет дедлайн (~15 сек),
// после которого контекст отменяется и бот падает с "context deadline exceeded".
func (b *Bot) Start(_ context.Context) {
	if b.bot == nil {
		return
	}
	var ctx context.Context
	ctx, b.cancel = context.WithCancel(context.Background())

	// Регистрируем меню команд в Telegram
	b.setCommands(ctx)

	if b.webhookURL != "" {
		// Webhook mode: регистрируем вебхук в Telegram
		ok, err := b.bot.SetWebhook(ctx, &bot.SetWebhookParams{
			URL: b.webhookURL,
		})
		if err != nil || !ok {
			b.log.Error("failed to set telegram webhook", slog.Any("err", err))
			return
		}
		b.log.Info("telegram bot started (webhook mode)", slog.String("url", b.webhookURL))
		// В режиме webhook bot.Start не нужен — обновления придут через HTTP handler
		return
	}

	// Long polling mode
	b.log.Info("telegram bot started (long polling mode)")
	b.bot.Start(ctx)
}

// setCommands регистрирует меню команд бота в Telegram.
func (b *Bot) setCommands(ctx context.Context) {
	commands := []models.BotCommand{
		{Command: "start", Description: "Начать работу"},
		{Command: "help", Description: "Список команд"},
		{Command: "remind", Description: "Создать напоминание"},
		{Command: "cancel", Description: "Отменить текущее действие"},
		{Command: "app", Description: "Открыть Mini App"},
	}
	ok, err := b.bot.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil || !ok {
		b.log.Error("failed to set bot commands", slog.Any("err", err))
		return
	}
	b.log.Info("telegram bot commands registered")
}

// WebhookHandler возвращает HTTP handler для приёма вебхуков от Telegram.
// Подключается к роутеру только в webhook-режиме.
func (b *Bot) WebhookHandler() http.HandlerFunc {
	if b.bot == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}
	}
	return b.bot.WebhookHandler()
}

// IsWebhookMode возвращает true если бот работает в режиме webhook.
func (b *Bot) IsWebhookMode() bool {
	return b.webhookURL != ""
}

// Stop останавливает бота.
func (b *Bot) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
}

// SendMessage отправляет текстовое сообщение в указанный чат.
func (b *Bot) SendMessage(ctx context.Context, chatID int64, text string) error {
	if b.bot == nil {
		b.log.Warn("telegram bot is not configured, skipping message", slog.Int64("chat_id", chatID))
		return nil
	}
	_, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}
	return nil
}

// SendMessageWithMarkup отправляет сообщение с inline-клавиатурой.
func (b *Bot) SendMessageWithMarkup(ctx context.Context, chatID int64, text string, markup models.ReplyMarkup) error {
	if b.bot == nil {
		b.log.Warn("telegram bot is not configured, skipping message", slog.Int64("chat_id", chatID))
		return nil
	}
	_, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: markup,
	})
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}
	return nil
}
