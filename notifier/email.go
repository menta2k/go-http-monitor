package notifier

import (
	"context"
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	host     string
	port     int
	from     string
	username string
	password string
}

func NewEmailSender(host string, port int, from, username, password string) *EmailSender {
	return &EmailSender{
		host:     host,
		port:     port,
		from:     from,
		username: username,
		password: password,
	}
}

func (s *EmailSender) Send(_ context.Context, to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, to, subject, body)

	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
