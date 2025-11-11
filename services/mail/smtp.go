package mail

import (
	"gopkg.in/gomail.v2"
)

type MailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var Mailer *MailerConfig

func InitMailer(cfg MailerConfig) {
	Mailer = &cfg
}

func NewDialer() *gomail.Dialer {
	return gomail.NewDialer(
		Mailer.Host,
		Mailer.Port,
		Mailer.Username,
		Mailer.Password,
	)
}
