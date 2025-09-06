package web

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	usersServices "github.com/vovanwin/template/app/internal/module/users/services"
	tpl "github.com/vovanwin/template/app/internal/module/web/templates"
	"github.com/vovanwin/template/app/internal/shared/middleware"
)

type ControllerDeps struct {
	fx.In
	Router         *chi.Mux
	SessionManager *scs.SessionManager
	UsersService   usersServices.UsersService
}

type WebController struct {
	sessionManager *scs.SessionManager
	usersService   usersServices.UsersService
}

func Controller(deps ControllerDeps) {
	controller := &WebController{
		sessionManager: deps.SessionManager,
		usersService:   deps.UsersService,
	}

	deps.Router.Route("/web", func(r chi.Router) {
		r.Get("/login", controller.GetLogin)
		r.Post("/login", controller.PostLogin)
		r.Get("/me", controller.GetMe)
		r.Post("/logout", controller.PostLogout)
	})
}

// GetLogin отображает страницу входа
func (c *WebController) GetLogin(w http.ResponseWriter, r *http.Request) {
	// Проверяем, авторизован ли пользователь
	if c.isAuthenticated(r.Context()) {
		http.Redirect(w, r, "/web/me", http.StatusSeeOther)
		return
	}

	csrfToken := middleware.GetCSRFToken(c.sessionManager, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.LoginPage("", csrfToken).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// PostLogin обрабатывает отправку формы входа
func (c *WebController) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.renderLoginForm(w, r, "Ошибка обработки формы")
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	// Валидация
	if email == "" {
		c.renderLoginForm(w, r, "Введите email")
		return
	}
	if password == "" {
		c.renderLoginForm(w, r, "Введите пароль")
		return
	}

	// Проверка через сервис пользователей
	if !c.validateCredentials(r.Context(), email, password) {
		c.renderLoginForm(w, r, "Неверный email или пароль")
		return
	}

	// Получаем данные пользователя для сохранения в сессии
	user, err := c.usersService.GetUserByEmail(r.Context(), email)
	if err != nil || user == nil {
		c.renderLoginForm(w, r, "Ошибка аутентификации")
		return
	}

	// Сохраняем данные в сессию
	c.sessionManager.Put(r.Context(), "user_id", user.ID)
	c.sessionManager.Put(r.Context(), "user_email", user.Email)
	c.sessionManager.Put(r.Context(), "last_activity", time.Now())
	c.sessionManager.Put(r.Context(), "session_ip", middleware.GetClientIP(r))

	// Проверяем, есть ли URL для редиректа после логина
	redirectURL := c.sessionManager.PopString(r.Context(), "redirect_after_login")
	if redirectURL == "" {
		redirectURL = "/web/me"
	}

	// HTMX редирект
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

// GetMe отображает личный кабинет
func (c *WebController) GetMe(w http.ResponseWriter, r *http.Request) {
	email := c.sessionManager.GetString(r.Context(), "user_email")
	if email == "" {
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.MePage(email).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// PostLogout выход из системы
func (c *WebController) PostLogout(w http.ResponseWriter, r *http.Request) {
	if err := c.sessionManager.Destroy(r.Context()); err != nil {
		// Логируем ошибку, но не показываем пользователю
		// logger.Error(r.Context(), "Failed to destroy session", zap.Error(err))
	}

	// HTMX редирект
	w.Header().Set("HX-Redirect", "/web/login")
	w.WriteHeader(http.StatusOK)
}

// renderLoginForm рендерит форму входа с ошибкой (для HTMX partial)
func (c *WebController) renderLoginForm(w http.ResponseWriter, r *http.Request, errMsg string) {
	csrfToken := middleware.GetCSRFToken(c.sessionManager, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Для HTMX запросов возвращаем только форму
	err := tpl.LoginForm(errMsg, csrfToken).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// isAuthenticated проверяет, авторизован ли пользователь
func (c *WebController) isAuthenticated(ctx context.Context) bool {
	email := c.sessionManager.GetString(ctx, "user_email")
	return email != ""
}

// validateCredentials проверяет учетные данные через сервис пользователей
func (c *WebController) validateCredentials(ctx context.Context, email, password string) bool {
	user, err := c.usersService.ValidateCredentials(ctx, email, password)
	if err != nil {
		// В реальном приложении логируйте ошибку
		return false
	}
	return user != nil
}
