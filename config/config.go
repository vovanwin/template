package config

import (
	"fmt"
	"github.com/spf13/viper"
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
		Greylog        Greylog  `mapstructure:"greyLog"`
		Version        string   `mapstructure:"version"`
	}

	Server struct {
		Address string `mapstructure:"address"`
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
