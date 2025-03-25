package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
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

func main() {
	sentry_dsn := os.Getenv("SENTRY_DSN")
	if sentry_dsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: sentry_dsn,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
	} else {
		log.Println("Sentry is disabled.")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/up", up)

	log.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Printf("Server failed to start: %v\n", err)
	}
}
