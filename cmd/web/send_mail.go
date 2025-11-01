package main

import (
	"log"
	"time"

	"githup.com/Mo277210/booking/internal/models"

	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {

	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)

		}
	}()
}

func sendMsg(m models.MailData) error {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	client, err := server.Connect()
	if err != nil {
		errorLog.Println("Error connecting to mail server:", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, "hello ,<strong>world</strong>")

	err = email.Send(client)
	if err != nil {
		log.Println("Error sending email:", err)
	} else {
		log.Println("email sent!")
	}
return err
}
