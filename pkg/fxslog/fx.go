package fxslog

import (
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	slogpretty "template/pkg/fxslog/slogPretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger() func() *slog.Logger {

	var log *slog.Logger

	env := viper.GetString("env")

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		graylogHost := viper.GetString("greyLog.host")
		graylogPort := viper.GetString("greyLog.port")
		gelfWriter, err := gelf.NewWriter(fmt.Sprintf("%s:%s", graylogHost, graylogPort))
		if err != nil {
			panic("Не может подключится к грейлогу")
		}
		log = slog.New(sloggraylog.Option{Level: slog.LevelInfo, Writer: gelfWriter}.NewGraylogHandler())
	}
	slog.SetDefault(log)

	return func() *slog.Logger {
		return log
	}

}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
