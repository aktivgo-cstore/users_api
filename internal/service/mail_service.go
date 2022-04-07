package service

import (
	"crypto/tls"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	host     = os.Getenv("SMTP_HOST")
	port     = os.Getenv("SMTP_PORT")
	user     = os.Getenv("SMTP_USER")
	password = os.Getenv("SMTP_PASSWORD")
)

type MailService struct {
	SmtpConn *mail.SMTPClient
}

func CreateMailService() (*MailService, error) {
	server := mail.NewSMTPClient()

	port, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.Encryption = mail.EncryptionSSL

	server.KeepAlive = true
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpConn, err := server.Connect()
	if err != nil {
		return nil, err
	}

	return &MailService{
		SmtpConn: smtpConn,
	}, nil
}

func (ms *MailService) SendActivationMail(to string, name string, link string) error {
	title := "Активация аккаунта на " + apiUrl
	content := fmt.Sprintf(
		`<div><h1>Уважаемый(ая) %s, для активации перейдите по ссылке</h1><a href="%s">%s</a></div>`,
		name, link, link,
	)

	email, err := ms.createEmail(to, title, content)
	if err != nil {
		return err
	}

	if err := ms.sendMail(email); err != nil {
		return err
	}

	log.Println("Activation mail send")

	return nil
}

func (ms *MailService) SendRestorePasswordMail(to string, name string, tempPassword string) error {
	title := "Восстановление пароля на " + apiUrl
	content := fmt.Sprintf(
		`<div><h1>Уважаемый(ая) %s, ваш временный пароль: %s</h1>
			<h2>Обязательно смените его сразу после входа</h2></div>`,
		name, tempPassword,
	)

	email, err := ms.createEmail(to, title, content)
	if err != nil {
		return err
	}

	if err := ms.sendMail(email); err != nil {
		return err
	}

	log.Println("Restore mail send")

	return nil
}

func (ms *MailService) createEmail(to string, title string, content string) (*mail.Email, error) {
	email := mail.NewMSG()
	email.SetFrom(user).AddTo(to).SetSubject(title)
	htmlBody := content
	email.SetBody(mail.TextHTML, htmlBody)

	if email.Error != nil {
		return nil, email.Error
	}

	return email, nil
}

func (ms *MailService) sendMail(email *mail.Email) error {
	err := email.Send(ms.SmtpConn)
	if err != nil {
		return err
	}

	return nil
}
