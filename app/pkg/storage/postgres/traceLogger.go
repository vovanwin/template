package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

type Logger struct {
	zapLogger *zap.Logger
}

func NewLoggerTracer(zapLogger *zap.Logger) *Logger {
	return &Logger{zapLogger: zapLogger.Named("pgx")}
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case tracelog.LogLevelTrace:
		fields = append(fields, zap.Any("PGX_LOG_LEVEL", level))
		l.zapLogger.Debug(msg, fields...)
	case tracelog.LogLevelDebug:
		l.zapLogger.Debug(msg, fields...)
	case tracelog.LogLevelInfo:
		l.zapLogger.Info(msg, fields...)
	case tracelog.LogLevelWarn:
		l.zapLogger.Warn(msg, fields...)
	case tracelog.LogLevelError:
		l.zapLogger.Error(msg, fields...)
	default:
		fields = append(fields, zap.Any("INVALID_PGX_LOG_LEVEL", level))
		l.zapLogger.Error(msg, fields...)
	}
}
