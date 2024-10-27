package database

import (
	"fmt"
	"time"

	"github.com/n9mi/db-email/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, logger *logrus.Logger, viperCfg *viper.Viper, numEmailSeed int) {
	defaultRecipient := viperCfg.GetString("EMAIL_DEFAULT_RECIPIENT")
	if defaultRecipient == "" {
		logger.Fatal("empty EMAIL_DEFAULT_RECIPIENT configuration")
	}

	tx := db.Begin()

	t := time.Now()
	emailFormats := []*entity.EmailFormat{
		{
			Subject:        fmt.Sprintf("Example subject - %s", "WITHOUT CC"),
			BodyFormat:     fmt.Sprintf("%s - %s - This is [VALUE_1], [VALUE_2], fortunately [VALUE_3]", t.Format("2006-01-02 15:04:05"), "WITHOUT CC"),
			NumCustomValue: 3,
		},
		{
			Subject:        fmt.Sprintf("Example subject - %s", "WITH CC"),
			BodyFormat:     fmt.Sprintf("%s - %s - This is [VALUE_1], [VALUE_2], fortunately [VALUE_3]", t.Format("2006-01-02 15:04:05"), "WITH CC"),
			NumCustomValue: 3,
		},
	}
	if err := tx.Create(&emailFormats).Error; err != nil {
		tx.Rollback()
		logger.Fatal(err)
	}

	var emailBroadcasts []*entity.EmailBroadcast
	for _, eF := range emailFormats {
		for i := 1; i <= numEmailSeed; i++ {
			customValue1 := fmt.Sprintf("Custom Value 1 - %d", i)
			customValue2 := fmt.Sprintf("Custom Value 2 - %d", i)
			customValue3 := fmt.Sprintf("Custom Value 3 - %d", i)

			emailBroadcasts = append(emailBroadcasts, &entity.EmailBroadcast{
				EmailFormatID: eF.ID,
				EmailDest:     defaultRecipient,
				Column1Value:  &customValue1,
				Column2Value:  &customValue2,
				Column3Value:  &customValue3,
				Status:        0,
			})
		}
	}
	if err := tx.Create(emailBroadcasts).Error; err != nil {
		tx.Rollback()
		logger.Fatal(err)
	}

	if err := tx.Commit().Error; err != nil {
		logger.Fatal(err.Error())
	}
}
