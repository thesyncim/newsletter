package main

import (
	"github.com/thesyncim/email"
	"log"
	"net/smtp"
)

const (
	smtpServer = "smtp.gmail.com"
	smtpPort   = "587"
	smtpEmail  = "thesyncim@gmail.com"
	password   = ""
)

func sendEmail(contactList []NewsletterListRecord, newsletter NewsletterRecord) {

	for _, list := range contactList {
		m := email.NewMessage(newsletter.Description, newsletter.Content)
		m.From = smtpEmail
		m.To = list.Emails
		m.BodyContentType = "text/html"

		err := email.Send(smtpServer+":"+smtpPort, smtp.PlainAuth("", smtpEmail, password, smtpServer), m)
		if err != nil {
			log.Println(err)
		}

	}

}
