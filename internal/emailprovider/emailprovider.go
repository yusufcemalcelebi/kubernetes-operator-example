package emailprovider

import (
	"context"
	"fmt"
)

type EmailProvider interface {
	SendEmail(ctx context.Context, senderEmail, recipientEmail, subject, body string) (string, error)
}

type ProviderType string

const (
	MailerSend ProviderType = "MailerSend"
	Mailgun    ProviderType = "Mailgun"
)

func NewEmailProvider(providerType ProviderType, apiKey, domain string) (EmailProvider, error) {
	switch providerType {
	case MailerSend:
		return NewMailerSendProvider(apiKey), nil
	case Mailgun:
		return NewMailgunProvider(apiKey, domain), nil
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", providerType)
	}
}
