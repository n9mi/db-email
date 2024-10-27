package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

func NewViper() (*viper.Viper, error) {
	viperCfg := viper.New()

	if _, err := os.Stat("./.env"); errors.Is(err, os.ErrNotExist) {
		viper.AutomaticEnv()
	} else {
		viperCfg.SetConfigType("env")
		viperCfg.SetConfigFile(".env")
		viperCfg.AddConfigPath("./")

		if err := viperCfg.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return viperCfg, nil
}
