package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var (
	configFile = "config.yml"
	configType = "yml"
)

type (
	Config struct {
		Debug          bool     `mapstructure:"debug"`
		ContextTimeout int      `mapstructure:"contextTimeout"`
		Server         Server   `mapstructure:"server"`
		Database       Database `mapstructure:"database"`
		Env            string   `mapstructure:"env"`
		Graylog        Greylog  `mapstructure:"greyLog"`
		Version        string   `mapstructure:"version"`
		JWT            JWT      `mapstructure:"jwt"`
		LogLevel       string   `mapstructure:"logLevel"`
	}

	Server struct {
		Address           string        `mapstructure:"address"`
		ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout" default:"60s"`
		GracefulTimeout   time.Duration `mapstructure:"grace_full_timeout" default:"8s"`
	}

	JWT struct {
		AccessTtl  time.Duration `mapstructure:"access_ttl"`
		RefreshTtl time.Duration `mapstructure:"refresh_ttl"`
		SighKey    string        `mapstructure:"signKey"`
	}

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		Scheme   string `mapstructure:"scheme"`
	}

	Nsq struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	Greylog struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
)

func NewConfig() Config {
	conf := &Config{}

	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable decode into config struct, %v", err)
	}
	return *conf
}

func InitConfig() {
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	viper.AutomaticEnv()

	if err != nil {
		fmt.Println(err.Error())
	}
}
