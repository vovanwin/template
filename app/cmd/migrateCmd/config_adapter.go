package migrateCmd

import "github.com/vovanwin/template/app/config"

// ConfigAdapter адаптирует config.Config для использования с platform migrator
type ConfigAdapter struct {
	cfg *config.Config
}

// NewConfigAdapter создает новый адаптер конфигурации
func NewConfigAdapter(cfg *config.Config) *ConfigAdapter {
	return &ConfigAdapter{cfg: cfg}
}

// GetHost возвращает хост базы данных
func (a *ConfigAdapter) GetHost() string {
	return a.cfg.PG.HostPG
}

// GetPort возвращает порт базы данных
func (a *ConfigAdapter) GetPort() string {
	return a.cfg.PG.PortPG
}

// GetUsername возвращает имя пользователя базы данных
func (a *ConfigAdapter) GetUsername() string {
	return a.cfg.PG.UserPG
}

// GetPassword возвращает пароль базы данных
func (a *ConfigAdapter) GetPassword() string {
	return a.cfg.PG.PasswordPG
}

// GetDatabase возвращает имя базы данных
func (a *ConfigAdapter) GetDatabase() string {
	return a.cfg.PG.DbNamePG
}

// GetSchema возвращает схему базы данных
func (a *ConfigAdapter) GetSchema() string {
	return a.cfg.PG.SchemePG
}

// GetSSLMode возвращает режим SSL
func (a *ConfigAdapter) GetSSLMode() string {
	return "disable"
}
