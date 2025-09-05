package config

import (
	"log/slog"
	"net"
	"os"
	"path"
	"time"

	"github.com/vovanwin/template/app/pkg/validator"

	"github.com/ilyakaznacheev/cleanenv"
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
		Host string `env-required:"true" yaml:"host" env:"APP_HOST" validate:"required"`
		Port string `env-required:"true" yaml:"port" env:"APP_PORT" validate:"required"`
		Env  string `env-required:"true" yaml:"env" env:"APP_ENV" validate:"required,oneof=local dev prod"`

		ContextTimeout    time.Duration `env-required:"true" yaml:"context_timeout" env:"APP_CONTEXT_TIMEOUT" validate:"required"`
		ReadHeaderTimeout time.Duration `env-required:"true" yaml:"read_header_timeout" env:"APP_READ_HEADER_TIMEOUT" default:"60s"`
		GracefulTimeout   time.Duration `env-required:"true" yaml:"grace_ful_timeout" env:"APP_GRACE_FUL_TIMEOUT" default:"8s"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"APP_LOG_LEVEL"  validate:"required,oneof=DEBUG INFO WARN ERROR"`
	}

	PG struct {
		HostPG     string `env-required:"true" yaml:"host"          env:"APP_HOST_PG" validate:"required"`
		PortPG     string `env-required:"true" yaml:"port"          env:"APP_PORT_PG" validate:"required"`
		UserPG     string `env-required:"true" yaml:"user"          env:"APP_USER_PG" validate:"required"`
		PasswordPG string `env-required:"true" yaml:"password"      env:"APP_PASSWORD_PG" validate:"required"`
		SchemePG   string `env-required:"true" yaml:"scheme"        env:"APP_SCHEME_PG" validate:"required"`
		DbNamePG   string `env-required:"true" yaml:"db"            env:"APP_DBNAME_PG" validate:"required"`
	}
	Rabbit struct {
		URI string `env-required:"true" yaml:"amqp_url" env:"APP_AMQP_URI" validate:"required"`
	}

	JWT struct {
		SignKey    string        `env-required:"true" yaml:"sign_key" env:"APP_SIGN_KEY" validate:"required"`
		TokenTTL   time.Duration `env-required:"true" yaml:"token_ttl" env:"APP_TOKEN_TTL" validate:"required"`
		RefreshTTL time.Duration `env-required:"true" yaml:"refresh_token_ttl" env:"APP_REFRESH_TOKEN_TTL" validate:"required"`
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
