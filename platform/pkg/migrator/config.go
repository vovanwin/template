package migrator

import (
	"fmt"
	"net"
)

// Config содержит конфигурацию для подключения к базе данных
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Schema   string
	SSLMode  string
}

// DatabaseConfig интерфейс для получения конфигурации БД
type DatabaseConfig interface {
	GetHost() string
	GetPort() string
	GetUsername() string
	GetPassword() string
	GetDatabase() string
	GetSchema() string
	GetSSLMode() string
}

// NewConfig создает новую конфигурацию из интерфейса
func NewConfig(cfg DatabaseConfig) *Config {
	return &Config{
		Host:     cfg.GetHost(),
		Port:     cfg.GetPort(),
		Username: cfg.GetUsername(),
		Password: cfg.GetPassword(),
		Database: cfg.GetDatabase(),
		Schema:   cfg.GetSchema(),
		SSLMode:  cfg.GetSSLMode(),
	}
}

// DSN возвращает строку подключения к PostgreSQL
func (c *Config) DSN() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s&search_path=%s",
		c.Username, c.Password,
		net.JoinHostPort(c.Host, c.Port),
		c.Database, sslMode, c.Schema,
	)
}
