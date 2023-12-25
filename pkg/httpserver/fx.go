package httpserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"log"
	"log/slog"
	"template/config"
	"template/pkg/validator"
)

var Module = fx.Provide(newModule)

func newModule(lifecycle fx.Lifecycle, config config.Config, logger *slog.Logger) *echo.Echo {

	handler := echo.New()

	handler.Use(middleware.RequestID())

	handler.Use(middleware.Recover())
	//handler.Use(slogecho.New(logger))

	handler.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:   true,
		LogUserAgent: true,
		LogRequestID: true,
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("X-Request-ID", v.RequestID),
					slog.String("Latency", (v.Latency).String()),

					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("X-Request-ID", v.RequestID),
					slog.String("UserAgent", v.UserAgent),
					slog.String("Latency", (v.Latency).String()),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	handler.Use(RequestIdTest)

	// setup handler validator as lib validator
	handler.Validator = validator.NewCustomValidator()

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			//logrus.Print("Start Http Server.")
			go func() {
				err := handler.Start(config.Server.Address)
				if err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			//logrus.Print("Stopping Http Server.")
			return handler.Shutdown(ctx)
		},
	})

	return handler
}
