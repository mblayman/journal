package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/mblayman/journal/model"
)

//go:embed go_templates
var templates embed.FS
var tmpl = template.Must(template.ParseFS(templates, "go_templates/index.html"))

func index(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func getWebhookAuth() (string, string) {
	webhookSecret := os.Getenv("ANYMAIL_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("ANYMAIL_WEBHOOK_SECRET not set.")
	}
	parts := strings.Split(webhookSecret, ":")
	if len(parts) != 2 {
		log.Fatalf("ANYMAIL_WEBHOOK_SECRET must be in format 'username:password', got: %s", webhookSecret)
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

	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentryDsn,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
		logger.Println("Sentry is enabled.")
	} else {
		logger.Println("Sentry is disabled.")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/up", up)
	username, password := getWebhookAuth()
	// TODO: implement the real processor
	processor := func(model.EmailContent) {}
	mux.HandleFunc("/webhook", webhookHandler(username, password, processor, logger))

	logger.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Printf("Server failed to start: %v\n", err)
	}
}
