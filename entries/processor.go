package entries

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mblayman/journal/model"
)

func MakeEmailContentProcessor(requiredToAddress string, logger *log.Logger) model.EmailContentProcessor {
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

		logger.Printf("Parsed date from Subject: %s", date.Format(time.RFC3339))
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
