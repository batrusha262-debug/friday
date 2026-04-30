package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func New(host, port, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	if m.host == "" || m.username == "" {
		fmt.Printf("[mailer] SMTP not configured — code for %s: %s\n", to, body)

		return nil
	}

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	addr := m.host + ":" + m.port

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body,
	)

	return smtp.SendMail(addr, auth, m.from, []string{to}, []byte(msg))
}
