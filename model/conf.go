package model

type Config struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	DatabaseDirectory  string
	MattEmailAddress   string
	ReplyToAddress     string
	RequiredToAddress  string
	SentryDSN          string
	WebhookSecret      string
	JournalUser        string
	JournalPassword    string
}
