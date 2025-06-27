package webhook

import (
	"bytes"
	"database/sql"
	"embed"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mblayman/journal/entries"
	"github.com/mblayman/journal/model"
	_ "modernc.org/sqlite"
)

//go:embed ses-payload-test.json
var testPayload embed.FS

func TestSESWebhookHandler(t *testing.T) {
	username := "testuser"
	password := "testpass"

	// Set up in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create required tables
	_, err = db.Exec(`
        CREATE TABLE entries_entry (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            body TEXT NOT NULL,
            "when" DATE NOT NULL,
            user_id BIGINT NOT NULL,
            CONSTRAINT user_per_day UNIQUE (user_id, "when")
        );
        CREATE TABLE entries_prompt (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            "when" DATE NOT NULL,
            message_id TEXT NOT NULL,
            user_id BIGINT NOT NULL
        )`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	// Insert a prompt to satisfy MakeEmailContentProcessor
	_, err = db.Exec(`INSERT INTO entries_prompt ("when", message_id, user_id) VALUES (?, ?, ?)`,
		"2025-06-26", "01000197af3bb44e-49442306-e083-404b-823f-6cd120a64512-000000", 1)
	if err != nil {
		t.Fatalf("Failed to insert prompt: %v", err)
	}

	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", log.LstdFlags)
	conf := model.Config{
		ReplyToAddress:    "journal.abcdef123@inbound.journeyinbox.com",
		RequiredToAddress: "journal@mail.journeyinbox.com",
	}
	processor := entries.MakeEmailContentProcessor(conf, db, logger)
	handler := SESWebhookHandler(username, password, processor, logger)

	// Load payload from embedded file
	payload, err := testPayload.ReadFile("ses-payload-test.json")
	if err != nil {
		t.Fatalf("Failed to read ses-payload-test.json: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/ses-webhook", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", recorder.Code)
	}
	if recorder.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got %q", recorder.Body.String())
	}

	// Verify database entry
	var body string
	var when time.Time // Use time.Time instead of string
	err = db.QueryRow(`SELECT body, "when" FROM entries_entry WHERE user_id = 1 AND "when" = '2025-06-26'`).Scan(&body, &when)
	if err != nil {
		t.Errorf("Failed to query entry: %v, logs: %s", err, logBuf.String())
	}
	expectedBody := `This should log to my SNS topic and send me a test message because I'm responding to the proper reply address.

Here is another line so that I can use this as test input for the webhook receiver and have a realistic multi-line test.`
	if body != expectedBody {
		t.Errorf("Expected body %q, got %q, logs: %s", expectedBody, body, logBuf.String())
	}
	expectedWhen := "2025-06-26"
	if when.Format("2006-01-02") != expectedWhen {
		t.Errorf("Expected when %q, got %q, logs: %s", expectedWhen, when.Format("2006-01-02"), logBuf.String())
	}

	// Check log for upsert
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, "Upserted entry for user 1 on 2025-06-26") {
		t.Errorf("Expected log to contain 'Upserted entry for user 1 on 2025-06-26', got %q", logOutput)
	}
}

func TestSESWebhookHandlerUnauthorized(t *testing.T) {
	username := "testuser"
	password := "testpass"
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", log.LstdFlags)
	processor := func(model.EmailContent) {}
	handler := SESWebhookHandler(username, password, processor, logger)

	req := httptest.NewRequest(http.MethodPost, "/ses-webhook", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Header().Get("WWW-Authenticate"), `Basic realm="Restricted"`) {
		t.Errorf("Expected WWW-Authenticate header, got %v", recorder.Header().Get("WWW-Authenticate"))
	}
}

func TestSESWebhookHandlerMethodNotAllowed(t *testing.T) {
	username := "testuser"
	password := "testpass"
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", log.LstdFlags)
	processor := func(model.EmailContent) {}
	handler := SESWebhookHandler(username, password, processor, logger)

	req := httptest.NewRequest(http.MethodGet, "/ses-webhook", nil)
	req.SetBasicAuth(username, password)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed, got %d", recorder.Code)
	}
}
