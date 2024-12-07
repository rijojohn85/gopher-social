package mailer

import (
	"bytes"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"html/template"
	"log"
	"time"
)

type MailTripMailer struct {
	fromEmail string
	dialer    *gomail.Dialer
}

func NewMailTripDialer(host, username, password, fromEmail string, port int) *MailTripMailer {
	dialer := gomail.NewDialer(host, port, username, password)
	return &MailTripMailer{
		fromEmail: fromEmail,
		dialer:    dialer,
	}
}

func (m *MailTripMailer) Send(
	templateFile, _, userEmail string,
	data any,
	_ bool,
) error {
	message := gomail.NewMessage()

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
	// Set email headers
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", userEmail)
	message.SetHeader("Subject", "Hello from the Mailtrap team")

	// Set email body
	message.SetBody("text/html", body.String())

	for i := 0; i < maxRetries; i++ {
		err := m.dialer.DialAndSend(message)
		if err != nil {
			log.Println("Error sending email: ", err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		} else {
			log.Println("Email sent")
			return nil
		}
	}
	return fmt.Errorf("failed to send email after %d retries", maxRetries)
	return nil
}
