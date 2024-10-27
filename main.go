package main

import (
	"github.com/n9mi/db-email/config"
)

func main() {
	viperCfg, err := config.NewViper()
	if err != nil {
		panic(err)
	}
	logger := config.NewLogrus(viperCfg)
	db := config.NewDatabase(viperCfg)

	appCfg := config.AppConfig{
		Config: viperCfg,
		Logger: logger,
		DB:     db,
	}
	numEmailSeed := viperCfg.GetInt("NUM_EMAIL_SEEDS")
	if numEmailSeed < 1 {
		logger.Fatal("NUM_EMAIL_SEEDS hasn't been configured yet")
	}

	services := config.Setup(&appCfg, numEmailSeed)
	services.EmailSenderRegService.SendAll()
	services.EmailSenderCcService.SendAll()
}
