package entries

import (
	"bytes"
	"database/sql"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/mblayman/journal/model"
	_ "modernc.org/sqlite"
)

func TestMakeEmailContentProcessor(t *testing.T) {
	tests := []struct {
		name         string
		requiredTo   string
		emailContent model.EmailContent
		expectedLog  string
		wantBody     string // Expected body after upsert
	}{
		{
			name:       "Matching To address with valid Subject and Text",
			requiredTo: "journal.abcdef1@email.journeyinbox.com",
			emailContent: model.EmailContent{
				To:      "JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>",
				Subject: "It's Wednesday, Mar. 26, 2025. How are you?",
				Text: `I got up this morning at 8:30 and brushed my teeth, then left to go to Cafe
Ibiza to meet with Jared. Lorem ipsum dolor sit amet, consectetur adipiscing
elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
anim id est laborum.

On Wed, Mar 26, 2025, 9:00 AM JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com> wrote:
> Reply to this prompt to update your journal.
>`,
			},
			expectedLog: "Upserted entry for user 1 on 2025-03-26",
			wantBody: `I got up this morning at 8:30 and brushed my teeth, then left to go to Cafe Ibiza to meet with Jared. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
		},
		{
			name:       "Non-matching To address with valid Subject and Text",
			requiredTo: "journal.abcdef1@email.journeyinbox.com",
			emailContent: model.EmailContent{
				To:      "JourneyInbox Journal <journal.xyz789@email.journeyinbox.com>",
				Subject: "It's Wednesday, Mar. 26, 2025. How are you?",
				Text:    "Some text\nOn Wed, Mar 26, 2025, 9:00 AM JourneyInbox Journal <journal.xyz789@email.journeyinbox.com> wrote:\n>",
			},
			expectedLog: "Invalid To address: journal.xyz789@email.journeyinbox.com",
			wantBody:    "", // No upsert, so expect empty body
		},
		{
			name:       "Invalid To format with valid Subject and Text",
			requiredTo: "journal.abcdef1@email.journeyinbox.com",
			emailContent: model.EmailContent{
				To:      "Invalid Format",
				Subject: "It's Wednesday, Mar. 26, 2025. How are you?",
				Text:    "Some text\nOn Wed, Mar 26, 2025, 9:00 AM Someone <someone@example.com> wrote:\n>",
			},
			expectedLog: "Invalid To address:",
			wantBody:    "", // No upsert, so expect empty body
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up in-memory SQLite database
			db, err := sql.Open("sqlite", ":memory:")
			if err != nil {
				t.Fatalf("Failed to open in-memory database: %v", err)
			}
			defer db.Close()

			// Create the entries_entry table
			createTable := `
                CREATE TABLE entries_entry (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    body TEXT NOT NULL,
                    "when" DATE NOT NULL,
                    user_id BIGINT NOT NULL,
                    CONSTRAINT user_per_day UNIQUE (user_id, "when")
                )`
			if _, err := db.Exec(createTable); err != nil {
				t.Fatalf("Failed to create table: %v", err)
			}

			// Set up logger
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", log.LstdFlags)

			// Create and run processor
			processor := MakeEmailContentProcessor(tt.requiredTo, db, logger)
			processor(tt.emailContent)

			// Check log output
			logOutput := strings.TrimSpace(logBuf.String())
			if !strings.Contains(logOutput, tt.expectedLog) {
				t.Errorf("Expected log output to contain %q, got %q", tt.expectedLog, logOutput)
			}

			// Verify upsert result
			var body string
			err = db.QueryRow(`SELECT body FROM entries_entry WHERE user_id = 1 AND "when" = ?`, "2025-03-26").Scan(&body)
			if tt.wantBody == "" {
				if err != sql.ErrNoRows {
					t.Errorf("Expected no row for user 1 on 2025-03-26, got body %q", body)
				}
			} else {
				if err != nil {
					t.Errorf("Failed to query body: %v", err)
				} else if body != tt.wantBody {
					t.Errorf("Expected body %q, got %q", tt.wantBody, body)
				}
			}
		})
	}
}

// TestParseSubjectDate tests the date parsing from the subject string.
func TestParseSubjectDate(t *testing.T) {
	tests := []struct {
		name        string
		subject     string
		wantDate    time.Time
		wantErr     bool
		errContains string
	}{
		{
			name:     "Valid subject with Re:",
			subject:  "Re: It's Wednesday, Mar. 26, 2025. How are you?",
			wantDate: time.Date(2025, time.March, 26, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Valid subject without Re:",
			subject:  "It's Wednesday, Mar. 26, 2025. How are you?",
			wantDate: time.Date(2025, time.March, 26, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:        "Invalid subject (no It's)",
			subject:     "Invalid Subject",
			wantErr:     true,
			errContains: "invalid Subject format: Invalid Subject",
		},
		{
			name:        "Invalid date",
			subject:     "It's Wednesday, Foo. 26, 2025. How are you?",
			wantErr:     true,
			errContains: "error parsing date from Subject \"Foo. 26, 2025\"",
		},
		{
			name:        "Missing suffix",
			subject:     "It's Wednesday, Mar. 26, 2025",
			wantErr:     true,
			errContains: "invalid Subject format: It's Wednesday, Mar. 26, 2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSubjectDate(tt.subject)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseSubjectDate() error = nil, want error containing %q", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("parseSubjectDate() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("parseSubjectDate() error = %v, want nil", err)
				return
			}
			if !got.Equal(tt.wantDate) {
				t.Errorf("parseSubjectDate() = %v, want %v", got, tt.wantDate)
			}
		})
	}
}
