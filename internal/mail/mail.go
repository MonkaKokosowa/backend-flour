package mail

import (
	"fmt"

	"github.com/MonkaKokosowa/backend-flour/internal/env"
	"github.com/mrz1836/go-sanitize"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type Message struct {
	From    string
	To      string
	Subject string
	Body    string
}

func SendMail(dialer *gomail.Dialer, message Message) {
	log.Info().Msg(fmt.Sprintf("Sending mail from: %s", message.From))
	dialer.DialAndSend(compose_message(message))
}

func GetDialer(environment env.Environment) *gomail.Dialer {
	d := gomail.NewDialer(
		environment.Dialer.Server,
		environment.Dialer.Port,
		environment.Dialer.Username,
		environment.Dialer.Password,
	)
	return d
}

func compose_message(message Message) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", message.From)
	m.SetHeader("To", message.To)
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/html", sanitize.HTML(message.Body))
	return m
}
