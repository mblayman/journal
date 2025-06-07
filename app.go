package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/mblayman/journal/entries"
	"github.com/mblayman/journal/model"
	"github.com/mblayman/journal/webhook"
	_ "modernc.org/sqlite"
)

//go:embed templates
var templates embed.FS
var tmpl = template.Must(template.ParseFS(templates, "templates/index.html"))

func index(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func getConfig() model.Config {
	config := model.Config{
		DatabaseDirectory: os.Getenv("DB_DIR"),
		MattEmailAddress:  os.Getenv("MATT_EMAIL_ADDRESS"),
		RequiredToAddress: os.Getenv("REQUIRED_TO_ADDRESS"),
		SendGridAPIKey:    os.Getenv("SENDGRID_API_KEY"),
		SentryDSN:         os.Getenv("SENTRY_DSN"),
		UseAWS:            os.Getenv("USE_AWS"),
		WebhookSecret:     os.Getenv("ANYMAIL_WEBHOOK_SECRET"),
	}
	return config
}

func getWebhookAuth(config model.Config) (string, string) {

	if config.WebhookSecret == "" {
		log.Fatal("ANYMAIL_WEBHOOK_SECRET not set.")
	}
	parts := strings.Split(config.WebhookSecret, ":")
	if len(parts) != 2 {
		log.Fatalf("ANYMAIL_WEBHOOK_SECRET must be in format 'username:password', got: %s", config.WebhookSecret)
	}

	username := parts[0]
	password := parts[1]
	return username, password
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	err := godotenv.Load()
	if err != nil {
		logger.Printf("Could not load .env file: %v (continuing without it)", err)
	}

	config := getConfig()

	if config.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDSN,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
		logger.Println("Sentry is enabled.")
	} else {
		logger.Println("Sentry is disabled.")
	}

	if config.RequiredToAddress == "" {
		log.Fatal("REQUIRED_TO_ADDRESS not set.")
	}

	if config.MattEmailAddress == "" {
		log.Fatal("MATT_EMAIL_ADDRESS not set.")
	}

	dbPath := "./db.sqlite3"
	if config.DatabaseDirectory != "" {
		dbPath = filepath.Join(config.DatabaseDirectory, "db.sqlite3")
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	logger.Printf("Opened database at %s.", dbPath)

	// Email gateway setup
	if config.SendGridAPIKey == "" {
		log.Fatal("SENDGRID_API_KEY not set.")
	}
	var emailGateway entries.EmailGateway
	emailGateway = entries.NewSendGridGateway(config.SendGridAPIKey)
	if config.UseAWS == "yes" {
		logger.Println("Using the AWS SES gateway.")
		emailGateway = entries.NewAmazonSESGateway(config)
	} else {
		logger.Println("Using the SendGrid gateway.")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/up", up)
	username, password := getWebhookAuth(config)
	processor := entries.MakeEmailContentProcessor(config.RequiredToAddress, db, logger)
	mux.HandleFunc("/webhook", webhook.WebhookHandler(username, password, processor, logger))

	entries.RunDailyEmailTask(db, emailGateway, config, logger)

	logger.Println("Server starting on port 8000...")
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		logger.Printf("Server failed to start: %v\n", err)
	}
}
