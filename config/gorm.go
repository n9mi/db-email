package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(viperCfg *viper.Viper) *gorm.DB {
	dbHost := viperCfg.GetString("DB_HOST")
	dbUser := viperCfg.GetString("DB_USER")
	dbPassword := viperCfg.GetString("DB_PASSWORD")
	dbName := viperCfg.GetString("DB_NAME")
	dbPort := viperCfg.GetInt("DB_PORT")
	dbSsl := viperCfg.GetString("DB_SSL_MODE")
	dbTimezone := viperCfg.GetString("DB_TIMEZONE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
		dbSsl,
		dbTimezone)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
