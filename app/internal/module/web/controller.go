package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	tpl "github.com/vovanwin/template/app/internal/module/web/templates"
)

type ControllerDeps struct {
	fx.In
	Router         *chi.Mux
	SessionManager *scs.SessionManager
}

type WebController struct {
	sessionManager *scs.SessionManager
}

func Controller(deps ControllerDeps) {
	controller := &WebController{
		sessionManager: deps.SessionManager,
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tpl.LoginPage("").Render(r.Context(), w)
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

	// Простая проверка (в реальном приложении проверяйте в БД)
	if !c.validateCredentials(email, password) {
		c.renderLoginForm(w, r, "Неверный email или пароль")
		return
	}

	// Сохраняем сессию
	c.sessionManager.Put(r.Context(), "user_email", email)

	// HTMX редирект
	w.Header().Set("HX-Redirect", "/web/me")
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

// renderLoginForm рендерит форму входа с ошибкой
func (c *WebController) renderLoginForm(w http.ResponseWriter, r *http.Request, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Временно используем полную страницу, пока не сгенерируем templ
	err := tpl.LoginPage(errMsg).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// isAuthenticated проверяет, авторизован ли пользователь
func (c *WebController) isAuthenticated(ctx context.Context) bool {
	email := c.sessionManager.GetString(ctx, "user_email")
	return email != ""
}

// validateCredentials проверяет учетные данные (заглушка)
func (c *WebController) validateCredentials(email, password string) bool {
	// В реальном приложении проверяйте в БД с хешированием пароля
	// Пока что простая проверка для демонстрации
	return email == "admin@example.com" && password == "password" ||
		email == "user@example.com" && password == "123456"
}
