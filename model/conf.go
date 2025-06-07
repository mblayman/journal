package model

type Config struct {
	DatabaseDirectory string
	MattEmailAddress  string
	RequiredToAddress string
	SendGridAPIKey    string
	SentryDSN         string
	WebhookSecret     string
}
