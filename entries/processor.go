package entries

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mblayman/journal/model"
)

func MakeEmailContentProcessor(requiredToAddress string, db *sql.DB, logger *log.Logger) model.EmailContentProcessor {
	return func(emailContent model.EmailContent) {
		if emailContent.ToAddress() != requiredToAddress {
			logger.Printf("Invalid To address: %s", emailContent.ToAddress())
			return
		}

		date, err := parseSubjectDate(emailContent.Subject)
		if err != nil {
			logger.Printf("Failed to parse date from Subject: %v", err)
			return
		}

		replyText := emailContent.Reply()
		userID := int64(1)                   // Fixed user_id value
		dateStr := date.Format("2006-01-02") // SQLite date format (YYYY-MM-DD)
		query := `
            INSERT INTO entries_entry (user_id, "when", body)
            VALUES (?, ?, ?)
            ON CONFLICT(user_id, "when") DO UPDATE SET body = excluded.body`
		_, err = db.Exec(query, userID, dateStr, replyText)
		if err != nil {
			logger.Printf("Failed to upsert entry for user %d on %s: %v", userID, dateStr, err)
			return
		}

		logger.Printf("Upserted entry for user %d on %s", userID, dateStr)
	}
}

// parseSubjectDate extracts and parses the date from a subject string.
// It expects a format like "It's [Weekday], [Month]. [Day], [Year]. How are you?"
// with an optional prefix (e.g., "Re: ") before "It's".
func parseSubjectDate(subject string) (time.Time, error) {
	const prefix = "It's "
	const suffix = ". How are you?"
	if !strings.Contains(subject, prefix) || !strings.HasSuffix(subject, suffix) {
		return time.Time{}, fmt.Errorf("invalid Subject format: %s", subject)
	}

	// Extract the part after "It's" and before ". How are you?"
	startIdx := strings.Index(subject, prefix)
	trimmed := subject[startIdx+len(prefix):]
	trimmed = strings.TrimSuffix(trimmed, suffix)
	parts := strings.SplitN(trimmed, ", ", 2) // "Wednesday, Mar. 26, 2025" → ["Wednesday", "Mar. 26, 2025"]
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid Subject date format: %s", trimmed)
	}

	dateStr := parts[1] // "Mar. 26, 2025"
	t, err := time.Parse("Jan. 2, 2006", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date from Subject %q: %v", dateStr, err)
	}

	return t, nil
}
