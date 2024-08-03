package notification

import (
	"fmt"
	"log"
	"supa_go_ltp_updater/config"
	"supa_go_ltp_updater/model"
	"time"

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
		m.SetBody("text/html", recipient.Body)

		if err := mail.Send(s, m); err != nil {
			log.Printf("Could not send email to %s: %v", recipient, err)
		} else {
			log.Printf("Email sent successfully to %s!", recipient)
		}
	}
}

func GetEntryEmailList(filteredStockSlice []model.Stock) EmailList {
	// Get the current date
	now := time.Now()

	// Format the date as "30 July 2024"
	formattedDate := now.Format("02 January 2006")

	emails := []string{"nishantkumar9995@gmail.com", "nishantdotk@gmail.com", "saraswatimahato1998@gmail.com"}
	emailList := EmailList{}

	// Build the HTML body with a list of stocks
	body := `<p>Below are entry stock for today:</p>
             <ul>`

	for _, stock := range filteredStockSlice {
		link := "https://finance.yahoo.com/chart/" + stock.Symbol + ".NS"
		body += `<li><a href="` + link + `">` + stock.Symbol + `</a> - ` + fmt.Sprintf("%f", stock.GetVolumeTimes()) + `</li>`
	}

	body += `</ul>`

	for _, email := range emails {
		emailList = emailList.PushEmail(Email{
			To:      email,
			Subject: fmt.Sprintf("Tops Picks 🚀 %s", formattedDate),
			Body:    body,
		})
	}

	return emailList
}
