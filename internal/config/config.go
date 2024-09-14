package config

import "github.com/spf13/viper"

const (
	configFile = "configs/config.yml"
)

func InitConfig() error {
	viper.SetConfigFile(configFile)

	return viper.ReadInConfig()
}
