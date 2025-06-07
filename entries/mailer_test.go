package entries

import (
	"bytes"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/mblayman/journal/model"
)

// MockEmailGateway is a testable implementation of EmailGateway.
type MockEmailGateway struct {
	sentPrompts []struct{ subject, body string }
}

func (m *MockEmailGateway) SendPrompt(toName, toEmail, fromName, fromEmail, subject, body string) (string, error) {
	m.sentPrompts = append(m.sentPrompts, struct{ subject, body string }{subject, body})
	return "mock-msg-" + time.Now().Format("20060102150405"), nil // Dummy message ID
}

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
				t.Fatalf("Failed to create entries_prompt table: %v", err)
			}

			// Create entries_entry table
			_, err = db.Exec(`
                CREATE TABLE entries_entry (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    "when" DATE NOT NULL,
                    body TEXT NOT NULL,
                    user_id BIGINT NOT NULL
                )`)
			if err != nil {
				t.Fatalf("Failed to create entries_entry table: %v", err)
			}

			// Pre-insert a sample entry into entries_entry for user_id=1
			_, err = db.Exec(`
    INSERT INTO entries_entry ("when", body, user_id) 
    VALUES (?, ?, ?)`,
				"2025-04-01", "Test line 1\nTest line 2\n\nTest line 3", 1)
			if err != nil {
				t.Fatalf("Failed to pre-insert entry into entries_entry: %v", err)
			}

			// Pre-insert prompts into entries_prompt
			for _, p := range tt.preInsert {
				_, err := db.Exec(`
                    INSERT INTO entries_prompt ("when", message_id, user_id) 
                    VALUES (?, ?, ?)`,
					p.when, p.messageID, p.userID)
				if err != nil {
					t.Fatalf("Failed to pre-insert prompt: %v", err)
				}
			}

			// Mock logger and gateway
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)
			mockGateway := &MockEmailGateway{}
			config := model.Config{
				MattEmailAddress:  "matt@example.com",
				RequiredToAddress: "test@example.com",
			}

			// Run the function
			SendDailyEmails(db, mockGateway, config, logger, tt.currentTime)

			// Check number of prompts sent
			logOutput := logBuf.String()
			sentCount := bytes.Count(logBuf.Bytes(), []byte("Sent prompt for"))
			if sentCount != tt.wantPrompts {
				t.Errorf("Expected %d prompts sent, got %d; log: %q", tt.wantPrompts, sentCount, logOutput)
			}
			if len(mockGateway.sentPrompts) != tt.wantPrompts {
				t.Errorf("Mock gateway expected %d prompts, got %d", tt.wantPrompts, len(mockGateway.sentPrompts))
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
