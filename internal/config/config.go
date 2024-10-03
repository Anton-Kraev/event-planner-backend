package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type EnvType string

const (
	EnvLocal EnvType = "local"
)

type (
	HttpServerConfig struct {
		Address     string        `mapstructure:"address"`
		Timeout     time.Duration `mapstructure:"timeout"`
		IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	}

	TimetableAPIConfig struct {
		Address string        `mapstructure:"address"`
		Timeout time.Duration `mapstructure:"timeout"`
	}

	RedisConfig struct {
		Address        string        `mapstructure:"address"`
		Password       string        `mapstructure:"password,omitempty"`
		DB             int           `mapstructure:"db,omitempty"`
		UserExpPeriod  time.Duration `mapstructure:"user_exp_period"`
		EventExpPeriod time.Duration `mapstructure:"event_exp_period"`
	}

	Config struct {
		Env          EnvType            `mapstructure:"env,omitempty"`
		HttpServer   HttpServerConfig   `mapstructure:"http_server"`
		TimetableAPI TimetableAPIConfig `mapstructure:"timetable_api"`
		Redis        RedisConfig        `mapstructure:"redis"`
	}
)

const (
	configFile = "configs/config.yml"
)

func MustInit() Config {
	var (
		errInitConfig = "error initializing config: %s"
		config        Config
	)

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf(errInitConfig, err.Error())
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf(errInitConfig, err.Error())
	}

	return config
}
