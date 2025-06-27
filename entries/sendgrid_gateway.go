package entries

import (
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridGateway implements EmailGateway using the SendGrid API.
type SendGridGateway struct {
	client *sendgrid.Client
}

// NewSendGridGateway creates a new SendGridGateway with the given API key.
func NewSendGridGateway(apiKey string) *SendGridGateway {
	return &SendGridGateway{
		client: sendgrid.NewSendClient(apiKey),
	}
}

// SendPrompt sends an HTML email prompt via SendGrid and returns the message ID.
func (g *SendGridGateway) SendPrompt(toName, toEmail, fromName, fromEmail, replyToAddress, subject, body string) (string, error) {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, "", body) // Empty plain text, HTML body

	response, err := g.client.Send(message)
	if err != nil {
		return "", err
	}

	if msgID, ok := response.Headers["X-Message-Id"]; ok && len(msgID) > 0 {
		return msgID[0], nil
	}
	return "unknown-" + time.Now().Format("2006-01-02T15:04:05Z"), nil // Fallback
}
