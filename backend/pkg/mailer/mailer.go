package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/resend/resend-go/v3"
)

const sendTimeout = 10 * time.Second

type Mailer struct {
	apiKey string
	from   string
}

func New(apiKey, from string) *Mailer {
	return &Mailer{
		apiKey: apiKey,
		from:   from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	if m.apiKey == "" {
		fmt.Printf("[mailer] Resend not configured — code for %s: %s\n", to, body)

		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	client := resend.NewClient(m.apiKey)
	_, err := client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	})
	if err != nil {
		return fmt.Errorf("resend send email: %w", err)
	}

	return nil
}
