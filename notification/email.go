package notification

import (
	"fmt"
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"

	"gopkg.in/mail.v2"
)

// Define the Email struct
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailList []Email

// ------------getters
func (e *Email) GetEmail() Email {
	return *e
}

func (e *EmailList) GetEmails() []Email {
	return *e
}

func (e *EmailList) PushEmail(email Email) []Email {
	*e = append(*e, email)
	return *e
}

// -------------------
// TODO: add context for sending emails for which task the mails are sent
func SendMails(emailList EmailList) {
	if len(emailList) == 0 {
		return
	}

	d := mail.NewDialer(config.AppConfig.Email.SMTPServer, config.AppConfig.Email.SMTPPort, config.AppConfig.Email.SMTPEmail, config.AppConfig.Email.SMTPPass)

	s, err := d.Dial()
	if err != nil {
		log.Fatalf("Could not establish connection: %v", err)
	}
	defer s.Close()

	for _, recipient := range emailList {
		m := mail.NewMessage()
		m.SetHeader("From", config.AppConfig.Email.SMTPEmail)
		m.SetHeader("To", recipient.To)
		m.SetHeader("Subject", recipient.Subject)
		m.SetBody("text/plain", recipient.Body)

		if err := mail.Send(s, m); err != nil {
			log.Printf("Could not send email to %s: %v", recipient, err)
		} else {
			log.Printf("Email sent successfully to %s!", recipient)
		}
	}
}

// TODO: remove the and create a better fucntion later
func FilteredStockToEmail(stocksData model.Stock, context string) Email {
	email := Email{
		To:      "nishantdotk@gmail.com",
		Subject: context,
		Body:    fmt.Sprintf("Got this stock: \n %s", stocksData.Symbol),
	}

	return email
}
