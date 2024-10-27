package service

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/n9mi/db-email/entity"
	"github.com/n9mi/db-email/model"
	"github.com/n9mi/db-email/repository"
	"github.com/n9mi/db-email/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type EmailSenderCcService struct {
	DB *gorm.DB
	EmailSenderService
	EmailFormatRepository    *repository.EmailFormatRepository
	EmailBroadcastRepository *repository.EmailBroadcastRepository
}

func NewEmailSenderCcService(viperCfg *viper.Viper, logger *logrus.Logger, db *gorm.DB,
	emailFormatRepository *repository.EmailFormatRepository,
	emailBroadcastRepository *repository.EmailBroadcastRepository) *EmailSenderCcService {
	return &EmailSenderCcService{
		EmailSenderService: EmailSenderService{
			ViperCfg: viperCfg,
			Logger:   logger,
		},
		DB:                       db,
		EmailFormatRepository:    emailFormatRepository,
		EmailBroadcastRepository: emailBroadcastRepository,
	}
}

func (s *EmailSenderCcService) getEmailBroadcastRec(done <-chan bool, emailFormat *entity.EmailFormat) <-chan *model.EmailModel {
	stream := make(chan *model.EmailModel)

	go func() {
		defer close(stream)

		rows, err := s.DB.Model(&entity.EmailBroadcast{}).
			Select("id, email_dest, column1_value, column2_value, column3_value, column4_value, column5_value").
			Where("email_format_id = ?", emailFormat.ID).Rows()
		if err != nil {
			s.Logger.Warnf("can't find email_broadcasts : %+v", err)
			return
		}

		defer rows.Close()
		for rows.Next() {
			select {
			case <-done:
				return
			default:
				var emailBroadcast entity.EmailBroadcast
				if err := s.DB.ScanRows(rows, &emailBroadcast); err != nil {
					s.Logger.Panicf("failed while scanning records on email_broadcast : %+v", err)
				}

				emailModel := &model.EmailModel{
					To:               emailBroadcast.EmailDest,
					Subject:          emailFormat.Subject,
					Body:             util.GenerateEmailBody(emailFormat.NumCustomValue, emailFormat.BodyFormat, &emailBroadcast),
					EmailBroadcastID: emailBroadcast.ID,
				}
				stream <- emailModel
			}
		}
	}()

	return stream
}

func (s *EmailSenderCcService) sendEmailAndGetRes(done <-chan bool, eMChan <-chan *model.EmailModel) <-chan *model.EmailSendingRes {
	stream := make(chan *model.EmailSendingRes)

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			case eM, ok := <-eMChan:
				if !ok {
					return
				}
				if eM != nil {
					errSend := s.Send(eM)
					t := time.Now()
					res := model.EmailSendingRes{
						Recipient:        eM.To,
						Error:            errSend,
						At:               t,
						EmailBroadcastID: eM.EmailBroadcastID,
					}

					stream <- &res
				}
			}
		}
	}()

	return stream
}

func (s *EmailSenderCcService) fanIn(done <-chan bool, fannedOutChan []<-chan *model.EmailSendingRes) <-chan *model.EmailSendingRes {
	var wg sync.WaitGroup
	stream := make(chan *model.EmailSendingRes)

	for _, c := range fannedOutChan {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for val := range c {
				select {
				case <-done:
					return
				case stream <- val:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		fmt.Println("waiting fanin stream close...")
		close(stream)
	}()

	return stream
}

func (s *EmailSenderCcService) updateResultToDB(done <-chan bool, numCpu int, stream <-chan *model.EmailSendingRes) *model.OverallResult {
	resultAtomic := new(model.OverallResultAtomic)
	var wg sync.WaitGroup

	for i := 0; i < numCpu; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for val := range stream {
				select {
				case <-done:
					return
				default:
					if val != nil {
						resultAtomic.NumTotal.Add(1)
						if val.Error != nil {
							resultAtomic.NumFailed.Add(1)
							s.EmailBroadcastRepository.UpdateStatusByID(s.DB, val.EmailBroadcastID, 1)
						} else {
							resultAtomic.NumSuccess.Add(1)
							s.EmailBroadcastRepository.UpdateStatusByID(s.DB, val.EmailBroadcastID, 2)
						}
					}
				}
			}
		}()
	}

	wg.Wait()
	return resultAtomic.GetOverallResultModel()
}

func (s *EmailSenderCcService) SendAll() {
	emailFormatID := 2 // 1 For without concurrency, 2 for with concurrency
	numCpu := runtime.NumCPU()
	start := time.Now()
	done := make(chan bool)
	defer close(done)

	var emailFormat entity.EmailFormat
	if err := s.EmailFormatRepository.FindByID(s.DB, &emailFormat, emailFormatID); err != nil {
		s.Logger.Panicf("can't find format with id %d: %+v", emailFormatID, err)
	}

	emailStream := s.getEmailBroadcastRec(done, &emailFormat)
	fannedOutCh := make([]<-chan *model.EmailSendingRes, numCpu)
	for i := range numCpu {
		fannedOutCh[i] = s.sendEmailAndGetRes(done, emailStream)
	}

	fannedInCh := s.fanIn(done, fannedOutCh)
	overallRes := s.updateResultToDB(done, numCpu, fannedInCh)
	overallRes.StartAt = start
	overallRes.EndAt = time.Now()
	s.Logger.Infof("successfully send %d emails and failed to send %d (total: %d) in %f seconds with concurrency",
		overallRes.NumSuccess, overallRes.NumFailed, overallRes.NumTotal, overallRes.EndAt.Sub(overallRes.StartAt).Seconds())
}
