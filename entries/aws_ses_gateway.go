package entries

import "github.com/mblayman/journal/model"

// AmazonSESGateway implements EmailGateway using the AWS SES API.
type AmazonSESGateway struct {
}

// NewAmazonSESGateway creates a new AmazonSESGateway.
func NewAmazonSESGateway(config model.Config) *AmazonSESGateway {
	return &AmazonSESGateway{}
}

// SendPrompt sends an HTML email prompt via AWS and returns the message ID.
func (g *AmazonSESGateway) SendPrompt(toName, toEmail, fromName, fromEmail, subject, body string) (string, error) {
	// TODO: implement this with the AWS Go SDK
	// TODO: Does SES provide a message ID like SendGrid did?
	return "", nil
}
