package ui

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vovanwin/platform/server"
	"github.com/vovanwin/template/internal/controller/ui/pages"
	"github.com/vovanwin/template/internal/pkg/events"
	"github.com/vovanwin/template/internal/pkg/jwt"
	"github.com/vovanwin/template/internal/pkg/timezone"
	"github.com/vovanwin/template/internal/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type UIController struct {
	bus             *events.Bus
	authService     *service.AuthService
	reminderService *service.ReminderService
	jwtService      jwt.JWTService
	log             *slog.Logger
}

type Deps struct {
	fx.In
	Bus             *events.Bus
	AuthService     *service.AuthService
	ReminderService *service.ReminderService
	JWTService      jwt.JWTService
	Log             *slog.Logger
}

func NewUIController(deps Deps) *UIController {
	return &UIController{
		bus:             deps.Bus,
		authService:     deps.AuthService,
		reminderService: deps.ReminderService,
		jwtService:      deps.JWTService,
		log:             deps.Log,
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
// Возвращает true если обработчик должен прекратить работу.
func (c *UIController) requireAuth(w http.ResponseWriter, r *http.Request) (userID string, stop bool) {
	uid, _, ok := c.extractUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return "", true
	}
	return uid, false
}

func (c *UIController) RegisterRoutes(ctx context.Context, mux *runtime.ServeMux, _ *grpc.Server) error {
	if err := c.RegisterEvents(mux); err != nil {
		return err
	}

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
	templ.Handler(pages.DashboardPage(csrf.Token(r))).ServeHTTP(w, r)
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

	templ.Handler(pages.ProfilePage(profile, csrf.Token(r))).ServeHTTP(w, r)
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

	reminders, err := c.reminderService.ListReminders(r.Context(), userID)
	if err != nil {
		c.log.Error("list reminders", slog.Any("err", err))
		http.Error(w, "Ошибка загрузки напоминаний", http.StatusInternalServerError)
		return
	}

	templ.Handler(pages.RemindersPage(reminders, csrf.Token(r))).ServeHTTP(w, r)
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
		Title       string `json:"title"`
		Description string `json:"description"`
		RemindAt    string `json:"remind_at"`
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

	// TODO: получить telegram_chat_id из профиля пользователя
	var telegramChatID int64

	_, err = c.reminderService.CreateReminder(r.Context(), userID, req.Title, req.Description, remindAt, telegramChatID)
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
	templ.Handler(pages.RemindersList(reminders)).ServeHTTP(w, r)
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
	templ.Handler(pages.RemindersList(reminders)).ServeHTTP(w, r)
}

// handleSettings — страница настроек (GET /settings).
func (c *UIController) handleSettings(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if _, stop := c.requireAuth(w, r); stop {
		return
	}
	templ.Handler(pages.SettingsPage(csrf.Token(r))).ServeHTTP(w, r)
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
