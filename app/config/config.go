package config

import (
	"fmt"
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
		fmt.Printf("Debug: Нет конфиг: %v\n", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		fmt.Printf("Error: Ошибка формирования env из переменных окружения: %v\n", err)
	}

	if err := validator.NewCustomValidator().Validate(cfg); err != nil {
		fmt.Printf("Error: [config] Отсуствуют обязательные конфиги: %v\n", err)
		os.Exit(1)
	}
	return cfg, nil
}

type (
	Config struct {
		Server   `yaml:"server"`
		Log      `yaml:"log"`
		PG       `yaml:"PG"`
		Rabbit   `yaml:"rabbit"`
		JWT      `yaml:"JWT"`
		Temporal `yaml:"temporal"`
	}

	Server struct {
		Host string `yaml:"host" env:"APP_HOST" validate:"required"`
		Port string `yaml:"port" env:"APP_PORT" validate:"required"`
		Env  string `yaml:"env" env:"APP_ENV" validate:"required,oneof=local dev prod"`

		ContextTimeout    time.Duration `yaml:"context_timeout" env:"APP_CONTEXT_TIMEOUT" validate:"required"`
		ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"APP_READ_HEADER_TIMEOUT" default:"60s"`
		GracefulTimeout   time.Duration `yaml:"grace_ful_timeout" env:"APP_GRACE_FUL_TIMEOUT" default:"8s"`
	}

	Log struct {
		Level string `yaml:"level" env:"APP_LOG_LEVEL"  validate:"required,oneof=DEBUG INFO WARN ERROR"`
	}

	PG struct {
		HostPG     string `yaml:"host" env:"APP_HOST_PG" validate:"required"`
		PortPG     string `yaml:"port" env:"APP_PORT_PG" validate:"required"`
		UserPG     string `yaml:"user" env:"APP_USER_PG" validate:"required"`
		PasswordPG string `yaml:"password" env:"APP_PASSWORD_PG" validate:"required"`
		SchemePG   string `yaml:"scheme" env:"APP_SCHEME_PG" validate:"required"`
		DbNamePG   string `yaml:"db" env:"APP_DBNAME_PG" validate:"required"`
	}
	Rabbit struct {
		URI string `yaml:"amqp_url" env:"APP_AMQP_URI" validate:"required"`
	}

	JWT struct {
		SignKey    string        `yaml:"sign_key" env:"APP_SIGN_KEY" validate:"required"`
		TokenTTL   time.Duration `yaml:"token_ttl" env:"APP_TOKEN_TTL" validate:"required"`
		RefreshTTL time.Duration `yaml:"refresh_token_ttl" env:"APP_REFRESH_TOKEN_TTL" validate:"required"`
	}

	Temporal struct {
		Host      string `yaml:"host" env:"APP_TEMPORAL_HOST" default:"localhost"`
		Port      int    `yaml:"port" env:"APP_TEMPORAL_PORT" default:"7233"`
		Namespace string `yaml:"namespace" env:"APP_TEMPORAL_NAMESPACE" default:"default"`
		TaskQueue string `yaml:"task_queue" env:"APP_TEMPORAL_TASK_QUEUE" default:"default-task-queue"`
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
