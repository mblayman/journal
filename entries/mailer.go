package entries

import (
	"database/sql"
	"log"
	"time"
)

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

// SendDailyEmails sends prompt emails to users, catching up on any missed days.
// Assumes there is always at least one existing prompt for user_id=1.
// now is the current time to use for determining the date range.
func SendDailyEmails(db *sql.DB, emailGateway EmailGateway, requiredToAddress, mattEmailAddress string, logger *log.Logger, now time.Time) {
	const userID = 1 // Fixed user ID
	loc, _ := time.LoadLocation("America/New_York")
	today := now.In(loc).Truncate(24 * time.Hour)
	logger.Printf("Today is %s", today.Format("2006-01-02"))

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
		subject := "It’s " + date.Weekday().String() + ", " + date.Format("Jan. 2, 2006") + ". How are you?"
		body := "<p>Reply to this prompt to update your journal.</p>"

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
