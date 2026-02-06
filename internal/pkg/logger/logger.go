package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Options параметры для создания логгера.
type Options struct {
	// Level уровень логирования: DEBUG, INFO, WARN, ERROR
	Level string
	// JSON если true — вывод в JSON (для прода), иначе цветной текст (для локальной разработки)
	JSON bool
}

// NewLogger создаёт slog.Logger и устанавливает его как глобальный (slog.Default).
//
// Локально (JSON=false): цветной вывод через tint, время в читаемом формате, source для быстрого перехода в IDE.
// Прод (JSON=true): структурированный JSON в stdout, без цветов, с source для трейсинга ошибок.
func NewLogger(opts Options) *slog.Logger {
	level := parseLevel(opts.Level)

	var handler slog.Handler
	if opts.JSON {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	}

	l := slog.New(handler)
	slog.SetDefault(l)

	return l
}

func parseLevel(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
