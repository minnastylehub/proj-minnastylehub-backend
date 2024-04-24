package sendemail

import (
	"fmt"
	"net/smtp"
)

func SendEmail(to, subject, body, emailFrom, emailPass string) error {
	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", emailFrom, emailPass, smtpHost)

	// Compose email
	email := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, emailFrom, []string{to}, []byte(email))
	if err != nil {
		return err
	}

	return nil
}