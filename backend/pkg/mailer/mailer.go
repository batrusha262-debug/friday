package mailer

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

const sendTimeout = 10 * time.Second

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
	if m.host == "" || m.username == "" || m.password == "" {
		fmt.Printf("[mailer] SMTP not configured — code for %s: %s\n", to, body)

		return nil
	}

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	addr := m.host + ":" + m.port

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body,
	)

	return sendMail(addr, m.host, auth, m.from, []string{to}, []byte(msg))
}

func sendMail(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	dialer := net.Dialer{Timeout: sendTimeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(sendTimeout)); err != nil {
		return fmt.Errorf("smtp deadline: %w", err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
		if err := c.StartTLS(config); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth: %w", err)
			}
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}

	for _, addr := range to {
		if err := c.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp rcpt %s: %w", addr, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp write: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return c.Quit()
}
