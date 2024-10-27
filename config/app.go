package config

import (
	"github.com/n9mi/db-email/database"
	"github.com/n9mi/db-email/repository"
	"github.com/n9mi/db-email/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	Config *viper.Viper
	Logger *logrus.Logger
	DB     *gorm.DB
}

type EmailSenderServices struct {
	EmailSenderRegService *service.EmailSenderRegService
	EmailSenderCcService  *service.EmailSenderCcService
}

func Setup(appCfg *AppConfig, numEmailSeed int) *EmailSenderServices {
	database.Drop(appCfg.DB, appCfg.Logger)
	database.Create(appCfg.DB, appCfg.Logger)
	database.Seed(appCfg.DB, appCfg.Logger, appCfg.Config, numEmailSeed)

	emailFormatRepository := repository.NewEmailFormatRepository()
	emailBroadcastRepository := repository.NewEmailBroadcastRepository()
	emailSenderRegService := service.NewEmailSenderRegService(appCfg.Config, appCfg.Logger, appCfg.DB,
		emailFormatRepository, emailBroadcastRepository)
	emailSenderCcService := service.NewEmailSenderCcService(appCfg.Config, appCfg.Logger, appCfg.DB,
		emailFormatRepository, emailBroadcastRepository)

	return &EmailSenderServices{
		EmailSenderRegService: emailSenderRegService,
		EmailSenderCcService:  emailSenderCcService,
	}
}
