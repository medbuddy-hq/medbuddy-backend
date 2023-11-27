package utility

import (
	"bytes"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	log "github.com/sirupsen/logrus"
	"html/template"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"time"
)

type Email struct {
	from       string
	subject    string
	to         string
	privateKey string
}

func NewEmail(from, subject, to, privateKey string) *Email {
	return &Email{
		from:       from,
		subject:    subject,
		to:         to,
		privateKey: privateKey,
	}
}

func (email *Email) SendReminderEmail(logger *log.Logger, data *model.MedicationForDosage) error {
	tpl, err := template.ParseFiles("utility/template/reminder.html")
	if err != nil {
		log.Error("Error parsing html template file, error: ", err)
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := tpl.Execute(buf, data); err != nil {
		log.Error("Error executing html template file, error: ", err)
		return err
	}

	mg := mailgun.NewMailgun(email.from, email.privateKey)

	subject := "Have you taken your meds?"
	sender := "MedBuddy HQ <" + constant.AppName + "@" + email.from + ">"

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, "", email.to)
	message.SetHtml(buf.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = mg.Send(ctx, message)
	if err != nil {
		logger.Errorf("Error sending email to '%v', error: %v", email.to, err.Error())
		return err
	}

	return nil
}
