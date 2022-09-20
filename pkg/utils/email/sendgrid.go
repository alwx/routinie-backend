package email

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridEmailSender struct {
	APIKey string
}

const (
	FromAddress = "noreply@routinie.com"
)

func (sender *SendgridEmailSender) Send(
	personalizations []*mail.Personalization,
	fromName string,
	subject string,
	body string,
) error {
	client := sendgrid.NewSendClient(sender.APIKey)

	content := mail.NewContent("text/html", body)

	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail(fromName, FromAddress))
	m.Subject = subject
	m.AddPersonalizations(personalizations...)
	m.AddContent(content)

	response, err := client.Send(m)
	if response != nil {
		fmt.Println(sender.APIKey)
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return err
}
