package mailgun

import (
	"context"

	mailgun "github.com/mailgun/mailgun-go/v3"
	"github.com/mattes/go-mail"
)

type Mailgun struct {
	mail.Provider

	client *mailgun.MailgunImpl
}

func New(domain, apiKey string) (*Mailgun, error) {
	return &Mailgun{
		client: mailgun.NewMailgun(domain, apiKey),
	}, nil
}

func (p *Mailgun) Send(m *mail.Mail) error {
	if err := m.Render(); err != nil {
		return err
	}

	msg := p.client.NewMessage(
		m.FromStr(),
		m.Subject,
		string(m.Text),
		m.ToStr())

	msg.SetHtml(string(m.HTML))

	// unsure if it makes sense for Mailgun.Send to allow passing in
	// own context.
	_, _, err := p.client.Send(context.Background(), msg)
	return err
}
