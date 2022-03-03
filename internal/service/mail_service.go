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

func SendActivationMail(to string, link string) error {
	server := mail.NewSMTPClient()

	port, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.Encryption = mail.EncryptionSSL

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(user).AddTo(to).SetSubject("Активация аккаунта на " + apiUrl)
	htmlBody := fmt.Sprintf(`<div><h1>Для активации перейдите по ссылке</h1><a href="%s">%s</a></div>`, link, link)
	email.SetBody(mail.TextHTML, htmlBody)

	if email.Error != nil {
		return err
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	log.Println("Mail send")

	return nil
}
