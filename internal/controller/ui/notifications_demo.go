package ui

import (
	"net/http"
	"time"

	"github.com/gorilla/csrf"

	"github.com/vovanwin/template/internal/controller/ui/pages"
	"github.com/vovanwin/template/internal/pkg/events"
)

func (c *UIController) handleNotificationsDemo(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, stop := c.requireAuth(w, r); stop {
		return
	}
	token := csrf.Token(r)
	c.Render(w, r, pages.NotificationsDemoPage(token), pages.NotificationsDemoContent(token))
}

func (c *UIController) handleNotificationsDemoSend(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userID, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	msgType := r.FormValue("type")
	message := r.FormValue("message")
	if msgType == "" || message == "" {
		http.Error(w, "type and message are required", http.StatusBadRequest)
		return
	}

	if err := c.bus.Publish(events.Event{
		UserID:  userID,
		Message: message,
		Type:    msgType,
	}); err != nil {
		c.log.Error("demo send notification", "err", err)
		http.Error(w, "Failed to send", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *UIController) handleNotificationsDemoBurst(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userID, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	messages := []struct {
		msg string
		typ string
	}{
		{"Подключение установлено", "success"},
		{"Загрузка данных...", "info"},
		{"Обнаружена аномалия в метриках", "warning"},
		{"Синхронизация завершена", "success"},
		{"Новое сообщение от администратора", "info"},
	}

	go func() {
		for _, m := range messages {
			_ = c.bus.Publish(events.Event{
				UserID:  userID,
				Message: m.msg,
				Type:    m.typ,
			})
			time.Sleep(800 * time.Millisecond)
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}
