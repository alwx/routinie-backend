package email

import (
	"bytes"
	"html/template"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"

	"habiko-go/pkg/models"
)

type EmailHandler struct {
	EmailSender models.EmailSender
}

type tokenData struct {
	Token string
}

func (handler *EmailHandler) SendUserEmailConfirmation(user models.User) {
	tmpl, err := template.ParseFiles("email_templates/sign_up_email_confirmation.html")
	if err != nil {
		return
	}

	data := tokenData{
		Token: *user.EmailConfirmationToken,
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return
	}

	personalizations := make([]*mail.Personalization, 0)
	personalization := mail.NewPersonalization()
	personalization.AddTos(&mail.Email{Address: user.Email})
	personalizations = append(personalizations, personalization)

	_ = handler.EmailSender.Send(personalizations, viper.GetString("email.from"), "Your confirmation code", tpl.String())
}

func (handler *EmailHandler) SendUserRemindPassword(user models.User) {
	tmpl, err := template.ParseFiles("email_templates/remind_password.html")
	if err != nil {
		return
	}

	data := tokenData{
		Token: *user.RemindPasswordToken,
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return
	}

	personalizations := make([]*mail.Personalization, 0)
	personalization := mail.NewPersonalization()
	personalization.AddTos(&mail.Email{Address: user.Email})
	personalizations = append(personalizations, personalization)

	_ = handler.EmailSender.Send(personalizations, viper.GetString("email.from"), "Your code to set up a new password", tpl.String())
}
