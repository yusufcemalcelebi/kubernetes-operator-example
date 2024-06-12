package emailprovider

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunProvider struct {
	client *mailgun.MailgunImpl
	domain string
}

func NewMailgunProvider(apiKey, domain string) EmailProvider {
	return &MailgunProvider{
		client: mailgun.NewMailgun(domain, apiKey),
		domain: domain,
	}
}

func (p *MailgunProvider) SendEmail(ctx context.Context, senderEmail, recipientEmail, subject, body string) (string, error) {
	m := p.client.NewMessage(senderEmail, subject, body, recipientEmail)

	_, id, err := p.client.Send(ctx, m)
	if err != nil {
		return "", err
	}

	return id, nil
}
