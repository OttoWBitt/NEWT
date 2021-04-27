package main

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func sendEmail(id string, email string) error {
	m := gomail.NewMessage()

	emailAcc := "newtrecovery@gmail.com"

	// Set E-Mail sender
	m.SetHeader("From", emailAcc)

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "NEWT Password recovery")

	baseRecoverURL := "http://localhost:3000/reset"
	var text string = "Por favor, use o codigo '%s' para poder resetar sua senha em %s"
	recover := fmt.Sprintf(text, id, baseRecoverURL)

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
