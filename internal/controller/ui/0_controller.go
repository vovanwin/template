package ui

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/internal/controller/ui/layouts"
	"github.com/vovanwin/template/internal/controller/ui/pages"
	"github.com/vovanwin/template/internal/pkg/events"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type UIController struct {
	bus *events.Bus
}

func NewUIController(bus *events.Bus) *UIController {
	return &UIController{
		bus: bus,
	}
}

func (c *UIController) RegisterRoutes(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
	// SSE события
	if err := c.RegisterEvents(mux); err != nil {
		return err
	}

	// Маршруты для страниц (через стандартный Handle)
	err := mux.HandlePath("GET", "/login", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		token := csrf.Token(r)
		slog.Debug("Generated CSRF Token", slog.String("path", "/login"), slog.Int("token_len", len(token)))
		templ.Handler(pages.LoginPage(token)).ServeHTTP(w, r)
	})
	if err != nil {
		return err
	}

	err = mux.HandlePath("GET", "/", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Главная страница
		token := csrf.Token(r)
		templ.Handler(layouts.Layout("Главная", templ.Raw("<h1 class='text-3xl font-bold text-gray-800'>Добро пожаловать в Template App!</h1>"), token)).ServeHTTP(w, r)
	})
	if err != nil {
		return err
	}

	return nil
}

// Module возвращает fx.Option для подключения UI контроллера.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewUIController),
		fx.Provide(
			fx.Annotate(
				func(srv *UIController) server.GatewayRegistrator {
					return srv.RegisterRoutes
				},
				fx.ResultTags(`group:"gateway_registrators"`),
			),
		),
	)
}
