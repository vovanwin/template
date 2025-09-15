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

		// Server component ports
		HTTPPort    string `yaml:"http_port" env:"APP_HTTP_PORT" default:"8080"`
		GRPCPort    string `yaml:"grpc_port" env:"APP_GRPC_PORT" default:"8081"`
		DebugPort   string `yaml:"debug_port" env:"APP_DEBUG_PORT" default:"8082"`
		SwaggerPort string `yaml:"swagger_port" env:"APP_SWAGGER_PORT" default:"8084"`

		// Feature flags
		EnableHTTP     bool `yaml:"enable_http" env:"APP_ENABLE_HTTP" default:"true"`
		EnableGRPC     bool `yaml:"enable_grpc" env:"APP_ENABLE_GRPC" default:"true"`
		EnableDebug    bool `yaml:"enable_debug" env:"APP_ENABLE_DEBUG" default:"true"`
		EnableSwagger  bool `yaml:"enable_swagger" env:"APP_ENABLE_SWAGGER" default:"true"`
		EnableTemporal bool `yaml:"enable_temporal" env:"APP_ENABLE_TEMPORAL" default:"true"`
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
