package service

import (
	"time"

	"github.com/n9mi/db-email/entity"
	"github.com/n9mi/db-email/model"
	"github.com/n9mi/db-email/repository"
	"github.com/n9mi/db-email/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type EmailSenderRegService struct {
	DB *gorm.DB
	EmailSenderService
	EmailFormatRepository    *repository.EmailFormatRepository
	EmailBroadcastRepository *repository.EmailBroadcastRepository
}

func NewEmailSenderRegService(viperCfg *viper.Viper, logger *logrus.Logger, db *gorm.DB,
	emailFormatRepository *repository.EmailFormatRepository,
	emailBroadcastRepository *repository.EmailBroadcastRepository) *EmailSenderRegService {
	return &EmailSenderRegService{
		EmailSenderService: EmailSenderService{
			ViperCfg: viperCfg,
			Logger:   logger,
		},
		DB:                       db,
		EmailFormatRepository:    emailFormatRepository,
		EmailBroadcastRepository: emailBroadcastRepository,
	}
}

func (s *EmailSenderRegService) SendAll() {
	emailFormatID := 1 // 1 For without concurrency, 2 for with concurrency
	overallRes := new(model.OverallResult)
	overallRes.StartAt = time.Now()

	var emailFormat entity.EmailFormat
	if err := s.EmailFormatRepository.FindByID(s.DB, &emailFormat, emailFormatID); err != nil {
		s.Logger.Panicf("can't find format with id %d: %+v", emailFormatID, err)
	}

	rows, err := s.DB.Model(&entity.EmailBroadcast{}).
		Select("id, email_dest, column1_value, column2_value, column3_value, column4_value, column5_value").
		Where("email_format_id = ?", emailFormatID).Rows()
	if err != nil {
		s.Logger.Panicf("can't find email_broadcasts : %+v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var emailBroadcast entity.EmailBroadcast
		if err := s.DB.ScanRows(rows, &emailBroadcast); err != nil {
			s.Logger.Panicf("failed while scanning records on email_broadcast : %+v", err)
		}
		emailModel := model.EmailModel{
			To:      emailBroadcast.EmailDest,
			Subject: emailFormat.Subject,
			Body:    util.GenerateEmailBody(emailFormat.NumCustomValue, emailFormat.BodyFormat, &emailBroadcast),
		}

		errSend := s.Send(&emailModel)
		overallRes.NumTotal += 1
		if errSend != nil {
			s.EmailBroadcastRepository.UpdateStatusByID(s.DB, emailBroadcast.ID, 1)
			overallRes.NumFailed += 1
		} else {
			s.EmailBroadcastRepository.UpdateStatusByID(s.DB, emailBroadcast.ID, 2)
			overallRes.NumSuccess += 1
		}
	}

	overallRes.EndAt = time.Now()
	s.Logger.Infof("successfully send %d emails and failed to send %d (total: %d) in %f seconds without concurrency",
		overallRes.NumSuccess, overallRes.NumFailed, overallRes.NumTotal, overallRes.EndAt.Sub(overallRes.StartAt).Seconds())
}
