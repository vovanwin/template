package ui

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vovanwin/template/internal/pkg/jwt"
)

func (c *UIController) RegisterEvents(mux *runtime.ServeMux) error {
	return mux.HandlePath("GET", "/api/v1/events", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		rc := http.NewResponseController(w)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // Для Nginx

		// Извлекаем ID пользователя из контекста (через наш AuthMiddleware)
		userID, ok := jwt.GetUserIDFromContext(r.Context())
		if !ok {
			userID = "anonymous"
		}

		eventsCh := c.bus.Subscribe(userID)
		defer c.bus.Unsubscribe(userID)

		// Начальное сообщение
		fmt.Fprintf(w, "data: <div class='hidden'>Connected</div>\n\n")
		rc.Flush()

		ticker := time.NewTicker(15 * time.Second) // Keep-alive
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case e, ok := <-eventsCh:
				if !ok {
					return
				}
				msg := fmt.Sprintf("<div class='bg-green-500 text-white p-3 rounded-lg shadow-xl mb-3 border-l-4 border-green-700 animate-fade-in'>%s</div>", e.Message)
				fmt.Fprintf(w, "data: %s\n\n", msg)
				rc.Flush()
			case <-ticker.C:
				fmt.Fprintf(w, ": keep-alive\n\n")
				rc.Flush()
			}
		}
	})
}
