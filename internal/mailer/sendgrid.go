package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"log"
	"time"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromEmail string, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(
	templateFile, username, email string,
	data any,
	isSandbox bool,
) error {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}
	message := mail.NewSingleEmail(
		from,
		subject.String(),
		to,
		"",
		body.String(),
	)
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})
	for i := 0; i < maxRetries; i++ {
		response, err := s.client.Send(message)
		if err != nil {
			log.Println("Error sending email: ", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		} else {
			log.Println("Email sent with response code ", response.StatusCode)
			return nil
		}
	}
	return fmt.Errorf("failed to send email after %d retries", maxRetries)
}
