package email

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func SendEmail(tokenRec string, email string) error {
	m := gomail.NewMessage()

	emailAcc := "newtrecovery@gmail.com"

	// Set E-Mail sender
	m.SetHeader("From", emailAcc)

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "NEWT Password recovery")

	baseRecoverURL := "http://newt.ottobittencourt.com:4200/reset/%s"
	recoverURL := fmt.Sprintf(baseRecoverURL, tokenRec)
	var text string = "Para recuperar sua senha, entre em %s"
	recover := fmt.Sprintf(text, recoverURL)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", recover)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, emailAcc, "zsvnjemrzrnnrmty")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
