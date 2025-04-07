package entries

import (
	"bytes"
	"database/sql"
	"log"
	"testing"
	"time"
)

func TestSendDailyEmails(t *testing.T) {
	tests := []struct {
		name        string
		currentTime time.Time // Mocked "now"
		preInsert   []struct {
			when, messageID string
			userID          int64
		}
		wantPrompts    int    // Expected number of prompts sent
		wantLastPrompt string // Expected last "when" value in DB
	}{
		{
			name:        "Last prompt yesterday, sends today",
			currentTime: time.Date(2025, 4, 6, 9, 0, 0, 0, time.FixedZone("EDT", -4*3600)),
			preInsert: []struct {
				when, messageID string
				userID          int64
			}{
				{when: "2025-04-05", messageID: "msg-2025-04-05", userID: 1},
			},
			wantPrompts:    1,
			wantLastPrompt: "2025-04-06",
		},
		{
			name:        "Missed two days, catches up",
			currentTime: time.Date(2025, 4, 8, 9, 0, 0, 0, time.FixedZone("EDT", -4*3600)),
			preInsert: []struct {
				when, messageID string
				userID          int64
			}{
				{when: "2025-04-05", messageID: "msg-2025-04-05", userID: 1},
			},
			wantPrompts:    3, // 04-06, 04-07, 04-08
			wantLastPrompt: "2025-04-08",
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

			// Create entries_prompt table
			_, err = db.Exec(`
                CREATE TABLE entries_prompt (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    "when" DATE NOT NULL,
                    message_id TEXT NOT NULL,
                    user_id BIGINT NOT NULL
                )`)
			if err != nil {
				t.Fatalf("Failed to create table: %v", err)
			}

			// Pre-insert prompts (always present)
			for _, p := range tt.preInsert {
				_, err := db.Exec(`INSERT INTO entries_prompt ("when", message_id, user_id) VALUES (?, ?, ?)`,
					p.when, p.messageID, p.userID)
				if err != nil {
					t.Fatalf("Failed to pre-insert prompt: %v", err)
				}
			}

			// Mock logger
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0) // No timestamp prefix

			// Run the function with mocked time
			SendDailyEmails(db, logger, tt.currentTime)

			// Check number of prompts sent (via log)
			logOutput := logBuf.String()
			sentCount := bytes.Count(logBuf.Bytes(), []byte("Sending prompt for"))
			if sentCount != tt.wantPrompts {
				t.Errorf("Expected %d prompts sent, got %d; log: %q", tt.wantPrompts, sentCount, logOutput)
			}

			// Verify the last prompt in the database
			var lastWhen string
			err = db.QueryRow(`SELECT strftime('%Y-%m-%d', "when") FROM entries_prompt WHERE user_id = 1 ORDER BY "when" DESC LIMIT 1`).Scan(&lastWhen)
			if err != nil {
				t.Errorf("Failed to query last prompt: %v", err)
			} else if lastWhen != tt.wantLastPrompt {
				t.Errorf("Expected last prompt date %q, got %q; log: %q", tt.wantLastPrompt, lastWhen, logOutput)
			}
		})
	}
}
