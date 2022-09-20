package models

import (
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSender interface {
	Send(personalizations []*mail.Personalization, fromName string, subject string, body string) error
}

type EmailHandler interface {
	SendUserEmailConfirmation(user User)
	SendUserRemindPassword(user User)
}
