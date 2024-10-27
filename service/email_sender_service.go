package service

import (
	"github.com/n9mi/db-email/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailSenderService struct {
	ViperCfg *viper.Viper
	Logger   *logrus.Logger
}

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587

func (s *EmailSenderService) Send(email *model.EmailModel) error {
	emailDefaultSender := s.ViperCfg.GetString("EMAIL_DEFAULT_SENDER")
	emailAppPwd := s.ViperCfg.GetString("EMAIL_APP_PASSWORD")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailDefaultSender)
	mailer.SetHeader("To", email.To)
	mailer.SetHeader("Subject", email.Subject)
	mailer.SetBody("text/html", email.Body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		emailDefaultSender,
		emailAppPwd,
	)
	return dialer.DialAndSend(mailer)
}
