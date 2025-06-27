package model

type Config struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	DatabaseDirectory  string
	MattEmailAddress   string
	ReplyToAddress     string
	RequiredToAddress  string
	SendGridAPIKey     string
	SentryDSN          string
	UseAWS             string
	WebhookSecret      string
}
