package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/ZaphCode/clean-arch/config"
)

type smtpEmailServiceImpl struct{}

func NewSmtpEmailService() EmailService {
	return new(smtpEmailServiceImpl)
}

func (s *smtpEmailServiceImpl) sendEmail(data EmailData) error {
	cfg := config.Get()

	subject := fmt.Sprintf("Subject: %s!\n", data.Subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	t, err := template.ParseFiles(
		// "./src/services/email/templates/" + data.Template,
		"./templates/" + data.Template,
	)

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, data.Data); err != nil {
		return err
	}

	body := buf.String()
	msg := []byte(subject + mime + body)

	sender := cfg.Smtp.Email
	password := cfg.Smtp.Password

	auth := smtp.PlainAuth("", sender, password, "smtp.gmail.com")

	recibers := []string{data.Email}

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		sender,
		recibers,
		[]byte(msg),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *smtpEmailServiceImpl) SendChangePasswordEmail(email, name, secretCode string) error {
	data := EmailData{
		Email:    email,
		Subject:  "Change Password Request",
		Template: "change_password.html",
		Data: map[string]interface{}{
			"Name":       name,
			"SecretCode": secretCode,
		},
	}
	return s.sendEmail(data)
}

func (s *smtpEmailServiceImpl) SendVerifyEmail(email, name, secretCode string) error {
	data := EmailData{
		Email:    email,
		Subject:  "Verify your Email",
		Template: "change_password.html",
		Data: map[string]interface{}{
			"Name":       name,
			"SecretCode": secretCode,
		},
	}
	return s.sendEmail(data)
}
