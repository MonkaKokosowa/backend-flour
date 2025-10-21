package mail

import (
	"fmt"

	"github.com/MonkaKokosowa/backend-flour/internal/env"
	"github.com/mrz1836/go-sanitize"
	"github.com/rs/zerolog/log"
	gomail "github.com/wneessen/go-mail"
)

type Address struct {
	Name  string
	Email string
}
type Message struct {
	From    Address
	To      string
	Subject string
	Body    string
	User    string
}

func LimitCharacters(input string, max int) string {
	if len(input) > max {
		return input[:max]
	}
	return input
}

func SendMail(client *gomail.Client, message Message) {
	log.Info().Msg(fmt.Sprintf("Sending mail from: %s", message.From))

	err := client.DialAndSend(compose_message(message))

	if err.Error() != "" {
		log.Error().Err(err)
	}
}

func GetClient(environment env.Environment) *gomail.Client {
	client, err := gomail.NewClient(
		environment.Dialer.Server,
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(environment.Dialer.Username),
		gomail.WithPassword(environment.Dialer.Password),
		gomail.WithPort(environment.Dialer.Port),
	)
	if err != nil {
		log.Error().Err(err)
	}
	return client
}

func compose_message(message Message) *gomail.Msg {
	m := gomail.NewMsg()
	m.SetAddrHeader(gomail.HeaderFrom, "Website Backend <"+message.User+">")
	m.To(message.To)
	m.Subject(message.Subject)
	m.SetBodyString(gomail.TypeTextPlain, sanitize.HTML("Message from "+message.From.Name+" "+message.From.Email+""+"\n\n"+message.Body))
	return m
}
