package entries

import (
	"bytes"
	"database/sql"
	"embed"
	"html/template"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

//go:embed templates
var templates embed.FS
var promptTmpl = template.Must(template.ParseFS(templates, "templates/prompt.html"))

// EmailGateway defines the interface for sending email prompts.
type EmailGateway interface {
	SendPrompt(toName, toEmail, fromName, fromEmail, subject, body string) (messageID string, err error)
}

// RunDailyEmailTask starts a goroutine that triggers SendDailyEmails daily at 9 AM Eastern Time.
func RunDailyEmailTask(db *sql.DB, emailGateway EmailGateway, requiredToAddress, mattEmailAddress string, logger *log.Logger) {
	go func() {
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			logger.Fatalf("Failed to load Eastern Time zone: %v", err)
		}

		for {
			now := time.Now().In(loc)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, loc)
			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}
			duration := nextRun.Sub(now)
			logger.Printf("Next email task scheduled for %s (in %v)", nextRun, duration)
			time.Sleep(duration)

			logger.Printf("Running daily email task at %s", time.Now().In(loc))
			SendDailyEmails(db, emailGateway, requiredToAddress, mattEmailAddress, logger, time.Now().In(loc))
		}
	}()
}

// linebreaks mimics Django's |linebreaks filter, converting newlines to HTML tags.
func linebreaks(text string) string {
	// Split into paragraphs by double newlines
	paragraphs := strings.Split(strings.TrimSpace(text), "\n\n")
	if len(paragraphs) == 0 {
		return ""
	}

	var htmlParts []string
	for _, para := range paragraphs {
		if para == "" {
			continue
		}
		// Replace single newlines within a paragraph with <br>
		lines := strings.Split(para, "\n")
		trimmedLines := make([]string, 0, len(lines))
		for _, line := range lines {
			if trimmed := strings.TrimSpace(line); trimmed != "" {
				trimmedLines = append(trimmedLines, trimmed)
			}
		}
		if len(trimmedLines) == 0 {
			continue
		}
		paraHTML := strings.Join(trimmedLines, "<br>")
		htmlParts = append(htmlParts, "<p>"+paraHTML+"</p>")
	}

	return strings.Join(htmlParts, "")
}

// createPromptBody generates the email body using a random entry for user_id=1.
func createPromptBody(db *sql.DB, today time.Time, logger *log.Logger) (string, error) {
	const userID = 1

	// Count entries for user_id=1
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM entries_entry WHERE user_id = ?`, userID).Scan(&count)
	if err != nil {
		return "", err
	}
	if count == 0 {
		logger.Printf("No entries found for user_id=%d", userID)
		return "<p>Reply to this prompt to update your journal.</p>", nil // Fallback
	}

	// Pick a random entry (OFFSET is 0-based)
	offset := rand.Intn(count)
	var whenStr, body string
	err = db.QueryRow(`
        SELECT strftime('%Y-%m-%d', "when"), body 
        FROM entries_entry 
        WHERE user_id = ? 
        ORDER BY "when" 
        LIMIT 1 OFFSET ?`, userID, offset).Scan(&whenStr, &body)
	if err != nil {
		return "", err
	}

	// Parse the entry's "when" date
	when, err := time.Parse("2006-01-02", whenStr)
	if err != nil {
		logger.Printf("Failed to parse entry date %q: %v", whenStr, err)
		return "", err
	}

	// Calculate delta using go-humanize
	delta := humanize.RelTime(when, today, "ago", "")

	// Convert body newlines to HTML
	bodyHTML := linebreaks(body)

	// Render the template
	data := struct {
		When  time.Time
		Delta string
		Body  template.HTML // Allows raw HTML from entry body
	}{
		When:  when,
		Delta: delta,
		Body:  template.HTML(bodyHTML),
	}
	var buf bytes.Buffer
	err = promptTmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SendDailyEmails sends prompt emails to users, catching up on any missed days.
// Assumes there is always at least one existing prompt for user_id=1.
// now is the current time to use for determining the date range.
func SendDailyEmails(db *sql.DB, emailGateway EmailGateway, requiredToAddress, mattEmailAddress string, logger *log.Logger, now time.Time) {
	const userID = 1 // Fixed user ID
	loc, _ := time.LoadLocation("America/New_York")
	today := now.In(loc).Truncate(24 * time.Hour)

	// Find the last prompt date for user_id=1
	var lastPromptDateStr string
	err := db.QueryRow(`
        SELECT strftime('%Y-%m-%d', "when") FROM entries_prompt 
        WHERE user_id = ? 
        ORDER BY "when" DESC 
        LIMIT 1`, userID).Scan(&lastPromptDateStr)
	if err != nil {
		logger.Printf("Failed to query last prompt: %v", err)
		return
	}
	lastPromptDate, err := time.Parse("2006-01-02", lastPromptDateStr)
	if err != nil {
		logger.Printf("Failed to parse last prompt date %q: %v", lastPromptDateStr, err)
		return
	}
	logger.Printf("Last prompt date: %s", lastPromptDate.Format("2006-01-02"))

	// Calculate missing days from last prompt to today
	startDate := lastPromptDate.Add(24 * time.Hour)
	logger.Printf("Start date for prompts: %s", startDate.Format("2006-01-02"))
	for date := startDate; !date.After(today); date = date.Add(24 * time.Hour) {
		// Construct email details
		toName := "Matt Layman"
		toEmail := mattEmailAddress
		fromName := "JourneyInbox Journal"
		fromEmail := requiredToAddress
		subject := "It's " + date.Weekday().String() + ", " + date.Format("Jan. 2, 2006") + ". How are you?"

		// Generate body with random entry
		body, err := createPromptBody(db, today, logger)
		if err != nil {
			logger.Printf("Failed to create prompt body for %s: %v", date.Format("2006-01-02"), err)
			continue
		}

		// Send email via gateway
		messageID, err := emailGateway.SendPrompt(toName, toEmail, fromName, fromEmail, subject, body)
		if err != nil {
			logger.Printf("Failed to send prompt for %s: %v", date.Format("2006-01-02"), err)
			continue
		}
		logger.Printf("Sent prompt for %s with message_id %s", date.Format("2006-01-02"), messageID)

		// Record the prompt in the database
		_, err = db.Exec(`
            INSERT INTO entries_prompt ("when", message_id, user_id) 
            VALUES (?, ?, ?)`,
			date.Format("2006-01-02"), messageID, userID)
		if err != nil {
			logger.Printf("Failed to insert prompt for %s: %v", date.Format("2006-01-02"), err)
			continue
		}
		logger.Printf("Recorded prompt for %s", date.Format("2006-01-02"))
	}
}
