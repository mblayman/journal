package entries

import (
	"log"
	"time"
)

// RunDailyEmailTask starts a goroutine that triggers SendDailyEmails daily at 9 AM Eastern Time.
func RunDailyEmailTask(logger *log.Logger) {
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
			SendDailyEmails(logger)
		}
	}()
}

// SendDailyEmails sends the daily emails to users.
// This is a placeholder for the actual email-sending implementation.
func SendDailyEmails(logger *log.Logger) {
	logger.Printf("Sending daily emails (placeholder)")
	// TODO: Implement email sending logic here
}
