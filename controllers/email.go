package controllers

import (
	"fmt"
	"io"
	"os"
	"github.com/go-gomail/gomail"
)

// // Function to send email with PDF attachment
// func SendEmailWithAttachment(to, subject, body, fileName, filePath string) error {
// 	SMTPemail := os.Getenv("Email")
// 	SMTPpass := os.Getenv("Password")
// 	smtpHost := "smtp.example.com"
// 	smtpPort := "587"

// 	// Compose email
// 	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", SMTPemail, to, subject, body)

// 	// Connect to the SMTP server
// 	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, smtpHost)
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, SMTPemail, []string{to}, []byte(msg))
// 	if err != nil {
// 		return err
// 	}

// 	// Read PDF file
// 	pdfData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return err
// 	}

// 	// Encode PDF attachment
// 	b64Data := base64.StdEncoding.EncodeToString(pdfData)

// 	// Compose email with attachment
// 	attachment := fmt.Sprintf("\r\nContent-Type: application/pdf\r\nContent-Disposition: attachment; filename=\"%s\"\r\nContent-Transfer-Encoding: base64\r\n\r\n%s", fileName, b64Data)
// 	msgWithAttachment := []byte(msg + attachment)

// 	// Send email with attachment
// 	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, SMTPemail, []string{to}, msgWithAttachment)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // SendEmail sends an email with an optional attachment
func SendEmail(msg, email, attachmentName string, attachmentData []byte) error {
	// SMTP server configuration
	senderEmail := os.Getenv("Email")
	senderPassword := os.Getenv("Password")

	// Compose email message
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Appointment Confirmation ")
	m.SetBody("text/plain", msg)

	// m.Attach(attachmentName, gomail.SetCopyFunc(func(w gomail.) error {
	// 	_, err := w.Write(attachmentData)
	// 	return err
	// }))

	// Add attachment
	m.Attach(attachmentName, gomail.SetCopyFunc(func(w io.Writer) error {
		// _, err := buf.WriteTo(w)
		_, err := w.Write(attachmentData)
		return err
	}))

	// Dial to SMTP server and send email
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
