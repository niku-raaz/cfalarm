package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go/v3"
	"github.com/sendgrid/sendgrid-go/v3/helpers/mail"
)

func SendReminderEmail(to string, content string) error {
	from := mail.NewEmail("cfalarm", "no-reply@cfalarm.com")
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, "Practice Reminder", toEmail, content, content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
