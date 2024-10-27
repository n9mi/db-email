package database

import (
	"github.com/n9mi/db-email/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Create(db *gorm.DB, logger *logrus.Logger) {
	if err := db.AutoMigrate(&entity.EmailFormat{}); err != nil {
		logger.Fatal(err)
	}

	if err := db.AutoMigrate(&entity.EmailBroadcast{}); err != nil {
		logger.Fatal(err)
	}
}

func Drop(db *gorm.DB, logger *logrus.Logger) {
	if err := db.Migrator().DropTable(&entity.EmailBroadcast{}); err != nil {
		logger.Fatal(err)
	}

	if err := db.Migrator().DropTable(&entity.EmailFormat{}); err != nil {
		logger.Fatal(err)
	}
}
