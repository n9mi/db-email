package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogrus(viperCfg *viper.Viper) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(viperCfg.GetInt("LOG_LEVEL")))
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	return log
}
