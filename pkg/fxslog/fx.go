package fxslog

import (
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"template/config"
	"template/pkg/fxslog/devslog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewLogger(c config.Config) *slog.Logger {
	var log *slog.Logger
	level := checkLevel(c)
	env := c.Env

	switch env {
	case envLocal:
		// new logger with options
		opts := &devslog.Options{
			HandlerOptions:    &slog.HandlerOptions{AddSource: true, Level: level},
			MaxSlicePrintSize: 4,
			SortKeys:          true,
			TimeFormat:        "[04:05]",
			NewLineAfterLog:   true,
			DebugColor:        devslog.Magenta,
		}
		log = slog.New(devslog.NewHandler(os.Stdout, opts))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	case envProd:
		graylogHost := viper.GetString("GRAYLOG_HOST")
		graylogPort := viper.GetString("GRAYLOG_PORT")
		gelfWriter, err := gelf.NewWriter(fmt.Sprintf("%s:%s", graylogHost, graylogPort))
		if err != nil {
			panic("Не может подключится к грейлогу")
		}
		log = slog.New(sloggraylog.Option{Level: level, Writer: gelfWriter}.NewGraylogHandler())
	}

	slog.SetDefault(log)

	return log

}

func checkLevel(c config.Config) slog.Level {
	if c.LogLevel == "debug" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}
