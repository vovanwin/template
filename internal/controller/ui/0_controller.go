package ui

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/internal/controller/ui/pages"
	"github.com/vovanwin/template/internal/pkg/centrifugo"
	"github.com/vovanwin/template/internal/pkg/events"
	"github.com/vovanwin/template/internal/pkg/jwt"
	"github.com/vovanwin/template/internal/pkg/timezone"
	"github.com/vovanwin/template/internal/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type UIController struct {
	bus              *events.Bus
	centrifugoClient *centrifugo.Client
	centrifugoURL    string
	authService      *service.AuthService
	reminderService  *service.ReminderService
	jwtService       jwt.JWTService
	log              *slog.Logger
}

type Deps struct {
	fx.In
	Bus              *events.Bus
	CentrifugoClient *centrifugo.Client
	CentrifugoURL    string `name:"centrifugo_url"`
	AuthService      *service.AuthService
	ReminderService  *service.ReminderService
	JWTService       jwt.JWTService
	Log              *slog.Logger
}

func NewUIController(deps Deps) *UIController {
	return &UIController{
		bus:              deps.Bus,
		centrifugoClient: deps.CentrifugoClient,
		centrifugoURL:    deps.CentrifugoURL,
		authService:      deps.AuthService,
		reminderService:  deps.ReminderService,
		jwtService:       deps.JWTService,
		log:              deps.Log,
	}
}

// extractUser достаёт userID и email из JWT-куки. Возвращает false если не авторизован.
func (c *UIController) extractUser(r *http.Request) (userID, email string, ok bool) {
	cookie, err := r.Cookie("access_token")
	if err != nil || cookie.Value == "" {
		return "", "", false
	}
	claims, err := c.jwtService.ValidateToken(cookie.Value)
	if err != nil || claims.TokenType != "access" {
		return "", "", false
	}
	return claims.UserID, claims.UserEmail, true
}

// requireAuth редиректит на /login если пользователь не авторизован.
// Для HTMX запросов возвращает 401 Unauthorized, что обрабатывается на клиенте.
func (c *UIController) requireAuth(w http.ResponseWriter, r *http.Request) (userID string, stop bool) {
	uid, _, ok := c.extractUser(r)
	if !ok {
		if r.Header.Get("HX-Request") == "true" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Сессия истекла"))
			return "", true
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return "", true
	}
	return uid, false
}

func (c *UIController) RegisterRoutes(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
	routes := []struct {
		method  string
		pattern string
		handler func(http.ResponseWriter, *http.Request, map[string]string)
	}{
		{"GET", "/", c.handleIndex},
		{"GET", "/login", c.handleLoginPage},
		{"POST", "/login", c.handleLoginSubmit},
		{"GET", "/logout", c.handleLogout},
		{"GET", "/dashboard", c.handleDashboard},
		{"GET", "/profile", c.handleProfile},
		{"GET", "/reminders", c.handleReminders},
		{"POST", "/reminders", c.handleCreateReminder},
		{"DELETE", "/reminders/{id}", c.handleDeleteReminder},
		{"GET", "/settings", c.handleSettings},
		{"GET", "/events-log", c.handleEventsLog},
		{"GET", "/notifications-demo", c.handleNotificationsDemo},
		{"POST", "/api/v1/notifications-demo/send", c.handleNotificationsDemoSend},
		{"POST", "/api/v1/notifications-demo/burst", c.handleNotificationsDemoBurst},
		{"GET", "/api/v1/centrifugo/token", c.handleCentrifugoToken},
	}

	for _, r := range routes {
		r := r // захват для замыкания
		if err := mux.HandlePath(r.method, r.pattern, r.handler); err != nil {
			return err
		}
	}

	return nil
}

// handleIndex — редирект в зависимости от состояния авторизации.
func (c *UIController) handleIndex(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, _, ok := c.extractUser(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// handleLoginPage — форма входа (GET /login).
func (c *UIController) handleLoginPage(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, _, ok := c.extractUser(r); ok {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	templ.Handler(pages.LoginPage(csrf.Token(r))).ServeHTTP(w, r)
}

// handleLoginSubmit — POST /login: прямой вызов сервиса, установка куки, HTMX-редирект.
func (c *UIController) handleLoginSubmit(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	result, err := c.authService.Login(r.Context(), req.Email, req.Password, r.RemoteAddr, r.UserAgent())
	if err != nil {
		c.log.Debug("login failed", slog.String("email", req.Email), slog.Any("err", err))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверный email или пароль"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 3600,
	})

	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

// handleLogout — сбрасывает куки и редиректит на /login.
func (c *UIController) handleLogout(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	for _, name := range []string{"access_token", "refresh_token"} {
		http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1})
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// handleDashboard — главная страница пользователя (GET /dashboard).
func (c *UIController) handleDashboard(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, stop := c.requireAuth(w, r); stop {
		return
	}
	token := csrf.Token(r)
	c.Render(w, r, pages.DashboardPage(token), pages.DashboardContent())
}

// handleProfile — страница профиля (GET /profile).
func (c *UIController) handleProfile(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userIDStr, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	profile, err := c.authService.GetProfile(r.Context(), userID)
	if err != nil {
		c.log.Error("get profile", slog.Any("err", err))
		http.Error(w, "Ошибка загрузки профиля", http.StatusInternalServerError)
		return
	}

	token := csrf.Token(r)
	c.Render(w, r, pages.ProfilePage(profile, token), pages.ProfileContent(profile, token))
}

// handleReminders — страница напоминаний (GET /reminders).
func (c *UIController) handleReminders(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userIDStr, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	const pageSize = 20
	paged, err := c.reminderService.ListRemindersPaged(r.Context(), userID, page, pageSize)
	if err != nil {
		c.log.Error("list reminders", slog.Any("err", err))
		http.Error(w, "Ошибка загрузки напоминаний", http.StatusInternalServerError)
		return
	}

	// Если это HTMX-запрос пагинации (не boosted навигация), возвращаем только таблицу
	if isHTMX(r) && r.Header.Get("HX-Boosted") != "true" {
		templ.Handler(pages.RemindersTablePaged(paged.Items, page, paged.TotalPages)).ServeHTTP(w, r)
		return
	}

	token := csrf.Token(r)
	c.Render(w, r,
		pages.RemindersPagePaged(paged.Items, token, page, paged.TotalPages),
		pages.RemindersContentPaged(paged.Items, token, page, paged.TotalPages),
	)
}

// handleCreateReminder — создание напоминания (POST /reminders).
func (c *UIController) handleCreateReminder(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	userIDStr, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	var req struct {
		Title                 string `json:"title"`
		Description           string `json:"description"`
		RemindAt              string `json:"remind_at"`
		RequireConfirmation   bool   `json:"require_confirmation"`
		RepeatIntervalMinutes int    `json:"repeat_interval_minutes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Парсим как локальное время пользователя → автоматически конвертируется в UTC при сохранении
	remindAt, err := timezone.FromUser("2006-01-02T15:04", req.RemindAt)
	if err != nil {
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	profile, err := c.authService.GetProfile(r.Context(), userID)
	if err != nil {
		c.log.Error("get profile for telegram_chat_id", slog.Any("err", err))
		http.Error(w, "Ошибка получения профиля", http.StatusInternalServerError)
		return
	}

	_, err = c.reminderService.CreateReminder(r.Context(), userID, req.Title, req.Description, remindAt, profile.TelegramChatID, req.RequireConfirmation, req.RepeatIntervalMinutes)
	if err != nil {
		c.log.Error("create reminder", slog.Any("err", err))
		http.Error(w, "Ошибка создания напоминания", http.StatusInternalServerError)
		return
	}

	// Возвращаем обновлённый список
	reminders, err := c.reminderService.ListReminders(r.Context(), userID)
	if err != nil {
		c.log.Error("list reminders", slog.Any("err", err))
		http.Error(w, "Ошибка загрузки", http.StatusInternalServerError)
		return
	}
	templ.Handler(pages.RemindersTable(reminders)).ServeHTTP(w, r)
}

// handleDeleteReminder — удаление напоминания (DELETE /reminders/{id}).
func (c *UIController) handleDeleteReminder(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	userIDStr, stop := c.requireAuth(w, r)
	if stop {
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	reminderID, err := uuid.Parse(pathParams["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := c.reminderService.DeleteReminder(r.Context(), userID, reminderID); err != nil {
		c.log.Error("delete reminder", slog.Any("err", err))
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}

	// Возвращаем обновлённый список
	reminders, err := c.reminderService.ListReminders(r.Context(), userID)
	if err != nil {
		c.log.Error("list reminders", slog.Any("err", err))
		http.Error(w, "Ошибка загрузки", http.StatusInternalServerError)
		return
	}
	templ.Handler(pages.RemindersTable(reminders)).ServeHTTP(w, r)
}

// handleSettings — страница настроек (GET /settings).
func (c *UIController) handleSettings(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, stop := c.requireAuth(w, r); stop {
		return
	}
	token := csrf.Token(r)
	c.Render(w, r, pages.SettingsPage(token), pages.SettingsContent())
}

// handleEventsLog — пример использования универсальной таблицы (GET /events-log).
func (c *UIController) handleEventsLog(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, stop := c.requireAuth(w, r); stop {
		return
	}

	// Мокаем данные для демонстрации универсальной таблицы
	mockEvents := []map[string]any{
		{"id": 1, "type": "auth", "msg": "User login success", "time": time.Now().Add(-1 * time.Hour).Format(time.RFC822)},
		{"id": 2, "type": "reminder", "msg": "Reminder created", "time": time.Now().Add(-2 * time.Hour).Format(time.RFC822)},
		{"id": 3, "type": "system", "msg": "SSE connected", "time": time.Now().Add(-3 * time.Hour).Format(time.RFC822)},
		{"id": 4, "type": "auth", "msg": "Token refreshed", "time": time.Now().Add(-4 * time.Hour).Format(time.RFC822)},
	}

	token := csrf.Token(r)
	c.Render(w, r, pages.EventsLogPage(mockEvents, token), pages.EventsLogContent(mockEvents))
}

// Module возвращает fx.Option для подключения UI контроллера.
func Module() fx.Option {
	return fx.Module("ui",
		fx.Decorate(func(log *slog.Logger) *slog.Logger {
			return log.With("component", "ui")
		}),
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
