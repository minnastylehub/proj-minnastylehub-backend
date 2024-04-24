package sendemail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

func ThankyouEmail(to, subject, body, emailFrom, emailPass string) error {
	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication
	auth := smtp.PlainAuth("", emailFrom, emailPass, smtpHost)

	// Get image data as base64
	imageBase64, err := getImageBase64()
	if err != nil {
		return err
	}

	// Create a new MIME message
	msg := createMIMEMessage(emailFrom, to, subject, body, imageBase64)

	// Convert message to bytes
	var buffer bytes.Buffer
	msg.WriteTo(&buffer)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, emailFrom, []string{to}, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func createMIMEMessage(emailFrom, to, subject, body, imageBase64 string) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)

	// Write email headers
	fmt.Fprintf(buf, "From: %s\r\n", emailFrom)
	fmt.Fprintf(buf, "To: %s\r\n", to)
	fmt.Fprintf(buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(buf, "Content-Type: multipart/related; boundary=boundary\r\n\r\n")

	// Write email body
	fmt.Fprintf(buf, "--boundary\r\n")
	fmt.Fprintf(buf, "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	fmt.Fprintf(buf, "%s\r\n\r\n", body)

	// Write image attachment
	fmt.Fprintf(buf, "--boundary\r\n")
	fmt.Fprintf(buf, "Content-Type: image/png\r\n")
	fmt.Fprintf(buf, "Content-Disposition: inline\r\n")
	fmt.Fprintf(buf, "Content-ID: <thankyou-image>\r\n")
	fmt.Fprintf(buf, "Content-Transfer-Encoding: base64\r\n\r\n")
	fmt.Fprintf(buf, "%s\r\n", imageBase64)

	// Write boundary end
	fmt.Fprintf(buf, "--boundary--\r\n")

	return buf
}

func getImageBase64() (string, error) {
	// Path to the image file
	imagePath := "./assets/thankyou-image.png"

	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	// Encode image data to base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	return imageBase64, nil
}
