package sessions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/vovanwin/template/app/config"
)

// SessionStoreType определяет тип хранилища сессий
type SessionStoreType string

const (
	Memory   SessionStoreType = "memory"
	Postgres SessionStoreType = "postgres"
	Redis    SessionStoreType = "redis"
)

// SessionProvider интерфейс для создания session manager
type SessionProvider interface {
	CreateSessionManager(config *config.Config, db *pgxpool.Pool) (*scs.SessionManager, error)
}

// DefaultSessionProvider реализация провайдера сессий
type DefaultSessionProvider struct{}

// NewSessionProvider создает новый провайдер сессий
func NewSessionProvider() SessionProvider {
	return &DefaultSessionProvider{}
}

// CreateSessionManager создает session manager в зависимости от конфигурации
func (p *DefaultSessionProvider) CreateSessionManager(cfg *config.Config, db *pgxpool.Pool) (*scs.SessionManager, error) {
	sessionManager := scs.New()

	// Базовая конфигурация безопасности
	sessionManager.Cookie.Name = "app_session"
	sessionManager.Lifetime = cfg.Sessions.Lifetime
	sessionManager.Cookie.HttpOnly = true                 // Защита от XSS
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode // CSRF защита

	// Дополнительные настройки безопасности
	sessionManager.Cookie.Path = "/"
	sessionManager.IdleTimeout = 30 * time.Minute // Автоматическое истечение при бездействии
	sessionManager.ErrorFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		// Логируем ошибки сессий, но не раскрываем детали пользователю
		http.Error(w, "Session error", http.StatusInternalServerError)
	}

	if cfg.IsProduction() {
		sessionManager.Cookie.Secure = true                      // Только HTTPS
		sessionManager.Cookie.SameSite = http.SameSiteStrictMode // Более строгая CSRF защита
	}

	// Выбор типа хранилища
	storeType := SessionStoreType(cfg.Sessions.Store)

	switch storeType {
	case Postgres:
		if db == nil {
			return nil, fmt.Errorf("database connection required for postgres session store")
		}

		// Создаем PostgreSQL store через stdlib
		sqlDB := stdlib.OpenDBFromPool(db)
		store := postgresstore.NewWithCleanupInterval(sqlDB, 15*time.Minute)
		sessionManager.Store = store

	case Redis:
		// TODO: Реализовать Redis store когда понадобится
		return nil, fmt.Errorf("redis session store not implemented yet")

	case Memory:
		// Используем встроенное in-memory хранилище (по умолчанию)
		// Подходит только для разработки

	default:
		return nil, fmt.Errorf("unsupported session store type: %s", storeType)
	}

	return sessionManager, nil
}

// CleanupExpiredSessions очищает истекшие сессии (для PostgreSQL)
func CleanupExpiredSessions(db *pgxpool.Pool) error {
	ctx := context.Background()
	_, err := db.Exec(ctx, "SELECT cleanup_expired_sessions()")
	return err
}
