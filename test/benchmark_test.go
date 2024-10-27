package test

import (
	"testing"

	"github.com/n9mi/db-email/config"
)

var services *config.EmailSenderServices

func init() {
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
	numEmailSeed := 5
	services = config.Setup(&appCfg, numEmailSeed)
}

func BenchmarkWithoutConcurrency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		services.EmailSenderRegService.SendAll()
	}
}

func BenchmarkWithConcurrency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		services.EmailSenderCcService.SendAll()
	}
}
