package handlers

import (
	"fmt"
	"net/smtp"
)

var domain string = "yourdomain.com"

func SendMail(mail string, nickname string, token string) bool {
	from := "no-reply@" + domain
	password := "mailpassword"

	smtpHost := "smtp." + domain
	smtpPort := "587"

	msg := []byte(fmt.Sprintf("Subject: Email Verification\n\n"+
		"Please verify your email by clicking the link:\n"+
		"http://"+domain+"/verify?token=%s", token))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{mail}, msg)
	return err == nil
}
