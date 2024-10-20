package fxslog

import (
	"fmt"
	"log/slog"
	"os"

	"app/pkg/fxslog/devslog"
)

const (
	envLocal = "local"
	envTest  = "test"
	envDev   = "dev"
	envProd  = "prod"
)

//go:generate options-gen -out-filename=slog_options.gen.go -from-struct=Options
type Options struct {
	level string `default:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR" `
	env   string `default:"prod" validate:"required,oneof=local dev prod test"`
}

func NewLogger(opts Options) (*slog.Logger, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options error: %w", err)
	}
	var log *slog.Logger
	level, err := ParseLevel(opts.level)
	if err != nil {
		return nil, fmt.Errorf("parse level error: %w", err)
	}

	switch opts.env {
	case envLocal, envTest:
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
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	}

	slog.SetDefault(log)

	return log, nil
}

func ParseLevel(text string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(text))
	return level, err
}
