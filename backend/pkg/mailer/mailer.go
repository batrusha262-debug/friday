package mailer

import (
	"context"
	"fmt"
	"mime"
	"net/smtp"
	"strings"
	"time"

	"github.com/resend/resend-go/v3"
)

const sendTimeout = 10 * time.Second

type Mailer struct {
	resendAPIKey string
	resendFrom   string
	smtpConfig   SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func New(apiKey, from string, smtpConfig SMTPConfig) *Mailer {
	return &Mailer{
		resendAPIKey: apiKey,
		resendFrom:   from,
		smtpConfig:   smtpConfig,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	if m.resendAPIKey != "" {
		return m.sendResend(to, subject, body)
	}

	if m.smtpConfig.Host != "" && m.smtpConfig.Username != "" && m.smtpConfig.Password != "" {
		return m.sendSMTP(to, subject, body)
	}

	if m.resendAPIKey == "" {
		fmt.Printf("[mailer] Resend not configured — code for %s: %s\n", to, body)

		return nil
	}

	return nil
}

func (m *Mailer) sendResend(to, subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	client := resend.NewClient(m.resendAPIKey)
	_, err := client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    m.resendFrom,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	})
	if err != nil {
		return fmt.Errorf("resend send email: %w", err)
	}

	return nil
}

func (m *Mailer) sendSMTP(to, subject, body string) error {
	from := m.smtpConfig.From
	if from == "" {
		from = m.smtpConfig.Username
	}

	addr := m.smtpConfig.Host + ":" + m.smtpConfig.Port
	auth := smtp.PlainAuth("", m.smtpConfig.Username, m.smtpConfig.Password, m.smtpConfig.Host)
	message := strings.Join([]string{
		"From: " + from,
		"To: " + to,
		"Subject: " + mime.QEncoding.Encode("utf-8", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message)); err != nil {
		return fmt.Errorf("smtp send email: %w", err)
	}

	return nil
}
