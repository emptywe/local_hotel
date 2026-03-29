package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/model"
	mail "github.com/xhit/go-simple-mail/v2"
	"go.uber.org/zap"
)

func ListenForMail(app *config.AppConfig) {
	zap.S().Debug("Mail listener start")
	for message := range app.MailChan {
		fmt.Println(message)
		sendMessage(message)
	}
}

func sendMessage(m model.MailData) {

	var message string

	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = time.Second * 10
	server.SendTimeout = time.Second * 10

	client, err := server.Connect()
	if err != nil {
		zap.S().Errorf("can't connect to mail server: %v", err)
	}

	email := mail.NewMSG()
	email = email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	if m.Template == "" {
		message = m.Content
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			zap.S().Errorf("can't parse mail template: %v", err)
		}
		message = strings.Replace(string(data), "[%body%]", m.Content, 1)
	}

	email.SetBody(mail.TextHTML, message)

	err = email.Send(client)
	if err != nil {
		zap.S().Errorf("can't send mail: %v", err)
	} else {
		zap.S().Debug("Email sent...")
	}
}

func sendMessageGoogle(m model.MailData) {

	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.KeepAlive = false
	server.ConnectTimeout = time.Second * 10
	server.SendTimeout = time.Second * 10
	server.Encryption = mail.EncryptionTLS
	server.Username = "@gmail.com"
	server.Password = ""

	client, err := server.Connect()
	if err != nil {
		zap.S().Errorf("can't connect to mail server: %v", err)
	}

	email := mail.NewMSG()
	email = email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		zap.S().Errorf("can't send mail: %v", err)
	} else {
		zap.S().Debug("Email sent...")
	}
}
