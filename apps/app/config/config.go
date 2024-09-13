package config

import (
	"app/pkg/validator"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"net"
	"os"
	"path"
	"time"
)

var configPath = "config/config.yml"

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		slog.Debug("Нет конфиг :", "err", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		slog.Error("Ошибка формирования env из переменных окружения:", "err", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	if err := validator.NewCustomValidator().Validate(cfg); err != nil {
		slog.Error("[config] Отсуствуют обязательные конфиги:", "err", err)
		os.Exit(1)
	}
	return cfg, nil
}

type (
	Config struct {
		Server `yaml:"server"`
		Log    `yaml:"log"`
		PG     `yaml:"PG"`
		Rabbit `yaml:"rabbit"`
		JWT    `yaml:"JWT"`
	}

	Server struct {
		Host string `env-required:"true" yaml:"host" env:"HOST" validate:"required"`
		Port string `env-required:"true" yaml:"port" env:"PORT" validate:"required"`
		Env  string `env-required:"true" yaml:"env" env:"ENV" validate:"required,oneof=local dev prod"`

		ContextTimeout    time.Duration `env-required:"true" yaml:"context_timeout" env:"CONTEXT_TIMEOUT" validate:"required"`
		ReadHeaderTimeout time.Duration `env-required:"true" yaml:"read_header_timeout" env:"READ_HEADER_TIMEOUT" default:"60s"`
		GracefulTimeout   time.Duration `env-required:"true" yaml:"grace_ful_timeout" env:"GRACE_FUL_TIMEOUT" default:"8s"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"  validate:"required,oneof=DEBUG INFO WARN ERROR"`
	}

	PG struct {
		HostPG     string `env-required:"true" yaml:"host"          env:"HOST_PG" validate:"required"`
		PortPG     string `env-required:"true" yaml:"port"          env:"PORT_PG" validate:"required"`
		UserPG     string `env-required:"true" yaml:"user"          env:"USER_PG" validate:"required"`
		PasswordPG string `env-required:"true" yaml:"password"      env:"PASSWORD_PG" validate:"required"`
		SchemePG   string `env-required:"true" yaml:"scheme"        env:"SCHEME_PG" validate:"required"`
		DbNamePG   string `env-required:"true" yaml:"db"            env:"DBNAME_PG" validate:"required"`
	}
	Rabbit struct {
		URI string `env-required:"true" yaml:"amqp_url" env:"AMQP_URI" validate:"required"`
	}

	JWT struct {
		SignKey    string        `env-required:"true" yaml:"sign_key" env:"SIGN_KEY" validate:"required"`
		TokenTTL   time.Duration `env-required:"true" yaml:"token_ttl" env:"TOKEN_TTL" validate:"required"`
		RefreshTTL time.Duration `env-required:"true" yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL" validate:"required"`
	}
)

func (c Config) Address() string {
	return net.JoinHostPort(c.Server.Host, c.Server.Port)
}

func (c Config) IsProduction() bool {
	return c.Server.Env == "prod"
}

func (c Config) IsLocal() bool {
	return c.Server.Env == "local"
}

func (c Config) IsTest() bool {
	return c.Server.Env == "test"
}
