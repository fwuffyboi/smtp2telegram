package main

import (
	"net/smtp"
	"strconv"
)

func sendmail() {
	username, password, host, port := "username", "password", "192.168.0.80", 25
	from, to, body := "googlecuck@gmail.com", "admin@skrunkle.cloud", "Hello, this is a test email."

	smtpServer := smtp.PlainAuth("", username, password, host)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Hello there\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(host+":"+strconv.Itoa(port), smtpServer, from, []string{to}, msg)
	if err != nil {
		panic(err)
	}
}
