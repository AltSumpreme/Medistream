package mail

import (
	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
)

var AsyncClient *asynq.Client

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", Mailer.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := NewDialer()
	return d.DialAndSend(m)
}
