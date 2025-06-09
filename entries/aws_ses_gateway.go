package entries

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/mblayman/journal/model"
)

// AmazonSESGateway implements EmailGateway using the AWS SES API.
type AmazonSESGateway struct {
	client *ses.Client
}

// NewAmazonSESGateway creates a new AmazonSESGateway.
func NewAmazonSESGateway(conf model.Config) *AmazonSESGateway {
	// Load AWS configuration with credentials from model.Config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.AWSAccessKeyID,
			conf.AWSSecretAccessKey,
			"",
		)),
		config.WithRegion("us-east-1"), // Adjust region as needed
	)
	if err != nil {
		// In a production scenario, you might want to handle this error differently
		// For now, panic to ensure the app doesn't start with an invalid config
		panic("failed to load AWS configuration: " + err.Error())
	}

	return &AmazonSESGateway{
		client: ses.NewFromConfig(cfg),
	}
}

// SendPrompt sends an HTML email prompt via AWS SES and returns the message ID.
func (g *AmazonSESGateway) SendPrompt(toName, toEmail, fromName, fromEmail, subject, body string) (string, error) {
	// Format the destination email address
	destination := &types.Destination{
		ToAddresses: []string{toName + " <" + toEmail + ">"},
	}

	// Create the email content
	content := &types.Content{
		Data:    aws.String(subject),
		Charset: aws.String("UTF-8"),
	}

	// Create the HTML body
	htmlBody := &types.Content{
		Data:    aws.String(body),
		Charset: aws.String("UTF-8"),
	}

	// Create the message body
	message := &types.Message{
		Subject: content,
		Body: &types.Body{
			Html: htmlBody,
		},
	}

	// Create the input for SendEmail
	input := &ses.SendEmailInput{
		Source:      aws.String(fromName + " <" + fromEmail + ">"),
		Destination: destination,
		Message:     message,
	}

	// Send the email
	result, err := g.client.SendEmail(context.TODO(), input)
	if err != nil {
		return "", err
	}

	// Return the SES Message ID
	return aws.ToString(result.MessageId), nil
}
