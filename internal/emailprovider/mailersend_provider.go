package emailprovider

import (
	"context"
	"time"

	"github.com/mailersend/mailersend-go"
)

type MailerSendProvider struct {
	client *mailersend.Mailersend
}

func NewMailerSendProvider(apiKey string) EmailProvider {
	return &MailerSendProvider{
		client: mailersend.NewMailersend(apiKey),
	}
}

func (p *MailerSendProvider) SendEmail(ctx context.Context, senderEmail, recipientEmail, subject, body string) (string, error) {
	message := p.client.Email.NewMessage()

	message.SetFrom(mailersend.From{Name: "Your Name", Email: senderEmail})
	message.SetRecipients([]mailersend.Recipient{{Name: "Recipient", Email: recipientEmail}})
	message.SetSubject(subject)
	message.SetHTML(body)
	message.SetText(body)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := p.client.Email.Send(ctx, message)
	if err != nil {
		return "", err
	}

	return res.Header.Get("X-Message-Id"), nil
}
