package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

func NewLogger() *slog.Logger {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			AddSource:  true,
			TimeFormat: time.DateTime,
		}),
	)

	return logger
}
