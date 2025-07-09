package main

import (
	"bufio"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/mblayman/journal/entries"
	"github.com/mblayman/journal/model"
	"github.com/mblayman/journal/webhook"
	_ "modernc.org/sqlite"
)

//go:embed templates
var templates embed.FS
var tmpl = template.Must(template.ParseFS(templates, "templates/index.html", "templates/journal.html"))

func index(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

// JournalEntry holds data for a single journal entry.
type JournalEntry struct {
	When       string
	Paragraphs []string
}

// MonthEntries holds entries for a specific month.
type MonthEntries struct {
	Name    string
	Entries []JournalEntry
}

// YearEntries holds months and their entries for a specific year.
type YearEntries struct {
	Months []MonthEntries
}

// JournalData holds the data for the journal template.
type JournalData struct {
	Years   []string
	Entries map[string]YearEntries
}

func journalHandler(db *sql.DB, config model.Config, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Basic authentication
		username, password, ok := r.BasicAuth()
		if !ok || username != config.JournalUser || password != config.JournalPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Journal"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get the 'when' query parameter
		requestedYear := r.URL.Query().Get("when")
		var yearFilter *int
		if requestedYear != "" {
			yearInt, err := strconv.Atoi(requestedYear)
			if err == nil && yearInt >= 1900 && yearInt <= time.Now().Year() {
				yearFilter = &yearInt
			} else {
				logger.Printf("Invalid 'when' parameter: %s", requestedYear)
				// Proceed without filtering if the year is invalid
			}
		}

		// Query all entries for user_id=1, ordered by date descending
		query := `SELECT strftime('%Y-%m-%d', "when"), body 
                  FROM entries_entry 
                  WHERE user_id = 1`
		args := []interface{}{}
		if yearFilter != nil {
			query += ` AND strftime('%Y', "when") = ?`
			args = append(args, fmt.Sprintf("%04d", *yearFilter))
		}
		query += ` ORDER BY "when" DESC`
		rows, err := db.Query(query, args...)
		if err != nil {
			logger.Printf("Failed to query journal entries: %v", err)
			http.Error(w, "Failed to fetch entries", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Organize entries by year and month
		entriesByYear := make(map[string]YearEntries)
		var years []string
		yearSet := make(map[string]bool)

		for rows.Next() {
			var whenStr, body string
			if err := rows.Scan(&whenStr, &body); err != nil {
				logger.Printf("Failed to scan entry: %v", err)
				http.Error(w, "Failed to process entries", http.StatusInternalServerError)
				return
			}

			// Parse date to extract year and month
			date, err := time.Parse("2006-01-02", whenStr)
			if err != nil {
				logger.Printf("Failed to parse date %s: %v", whenStr, err)
				continue
			}
			year := date.Format("2006")
			month := date.Format("January")

			// Skip entries not matching the year filter, if set
			if yearFilter != nil && year != fmt.Sprintf("%04d", *yearFilter) {
				continue
			}

			// Split body into paragraphs
			paragraphs := strings.Split(strings.TrimSpace(body), "\n\n")
			trimmedParagraphs := make([]string, 0, len(paragraphs))
			for _, para := range paragraphs {
				if trimmed := strings.TrimSpace(para); trimmed != "" {
					trimmedParagraphs = append(trimmedParagraphs, trimmed)
				}
			}

			// Add to entriesByYear
			yearEntries, exists := entriesByYear[year]
			if !exists {
				yearEntries = YearEntries{Months: []MonthEntries{}}
				entriesByYear[year] = yearEntries
				if !yearSet[year] {
					years = append(years, year)
					yearSet[year] = true
				}
			}

			// Find or create month
			var monthEntries *MonthEntries
			for i, m := range yearEntries.Months {
				if m.Name == month {
					monthEntries = &yearEntries.Months[i]
					break
				}
			}
			if monthEntries == nil {
				yearEntries.Months = append(yearEntries.Months, MonthEntries{Name: month, Entries: []JournalEntry{}})
				monthEntries = &yearEntries.Months[len(yearEntries.Months)-1]
			}

			// Add entry to month
			monthEntries.Entries = append(monthEntries.Entries, JournalEntry{When: whenStr, Paragraphs: trimmedParagraphs})
			entriesByYear[year] = yearEntries
		}
		if err := rows.Err(); err != nil {
			logger.Printf("Error iterating rows: %v", err)
			http.Error(w, "Failed to process entries", http.StatusInternalServerError)
			return
		}

		// Sort months within each year (in reverse chronological order)
		for year, yearEntries := range entriesByYear {
			sortedMonths := make([]MonthEntries, len(yearEntries.Months))
			copy(sortedMonths, yearEntries.Months)
			for i, month := range sortedMonths {
				monthTime, _ := time.Parse("January", month.Name)
				sortedMonths[i].Name = monthTime.Format("January") // Ensure consistent formatting
			}
			// Sort months by parsing month name to time
			for i := 0; i < len(sortedMonths)-1; i++ {
				for j := i + 1; j < len(sortedMonths); j++ {
					monthI, _ := time.Parse("January", sortedMonths[i].Name)
					monthJ, _ := time.Parse("January", sortedMonths[j].Name)
					if monthJ.After(monthI) {
						sortedMonths[i], sortedMonths[j] = sortedMonths[j], sortedMonths[i]
					}
				}
			}
			entriesByYear[year] = YearEntries{Months: sortedMonths}
		}

		// Sort years in descending order
		for i := 0; i < len(years)-1; i++ {
			for j := i + 1; j < len(years); j++ {
				if years[j] > years[i] {
					years[i], years[j] = years[j], years[i]
				}
			}
		}

		// Prepare data for template
		data := JournalData{
			Years:   years,
			Entries: entriesByYear,
		}

		// Set Content-Type before rendering
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.ExecuteTemplate(w, "journal.html", data)
		if err != nil {
			logger.Printf("Failed to render journal template: %v", err)
			// Avoid calling http.Error since response may have started
			return
		}
	}
}

func getConfig() model.Config {
	config := model.Config{
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DatabaseDirectory:  os.Getenv("DB_DIR"),
		MattEmailAddress:   os.Getenv("MATT_EMAIL_ADDRESS"),
		ReplyToAddress:     os.Getenv("REPLY_TO_ADDRESS"),
		RequiredToAddress:  os.Getenv("REQUIRED_TO_ADDRESS"),
		SentryDSN:          os.Getenv("SENTRY_DSN"),
		WebhookSecret:      os.Getenv("ANYMAIL_WEBHOOK_SECRET"),
		JournalUser:        os.Getenv("JOURNAL_USER"),
		JournalPassword:    os.Getenv("JOURNAL_PASSWORD"),
	}
	return config
}

func getWebhookAuth(config model.Config) (string, string) {
	if config.WebhookSecret == "" {
		log.Fatal("ANYMAIL_WEBHOOK_SECRET not set.")
	}
	parts := strings.Split(config.WebhookSecret, ":")
	if len(parts) != 2 {
		log.Fatalf("ANYMAIL_WEBHOOK_SECRET must be in format 'username:password', got: %s", config.WebhookSecret)
	}
	username := parts[0]
	password := parts[1]
	return username, password
}

// fixParagraphs removes newlines within paragraphs, preserving double newlines between them.
func fixParagraphs(text string) string {
	// Split into paragraphs by double newlines
	paragraphs := strings.Split(strings.TrimSpace(text), "\n\n")
	var fixedParagraphs []string
	for _, para := range paragraphs {
		if para == "" {
			continue
		}
		// Join lines within a paragraph, replacing single newlines with spaces
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
		fixedPara := strings.Join(trimmedLines, " ")
		fixedParagraphs = append(fixedParagraphs, fixedPara)
	}
	return strings.Join(fixedParagraphs, "\n\n")
}

// fixEntryForDate processes an entry for the given date: displays it, shows corrected version, and updates if confirmed.
func fixEntryForDate(db *sql.DB, date time.Time, logger *log.Logger) error {
	const userID = 1
	dateStr := date.Format("2006-01-02")
	var body string
	err := db.QueryRow(`
        SELECT body 
        FROM entries_entry 
        WHERE user_id = ? AND "when" = ?`, userID, dateStr).Scan(&body)
	if err == sql.ErrNoRows {
		logger.Printf("No entry found for user %d on %s", userID, dateStr)
		return fmt.Errorf("no entry found for %s", dateStr)
	}
	if err != nil {
		logger.Printf("Failed to query entry for user %d on %s: %v", userID, dateStr, err)
		return err
	}

	fmt.Printf("\nCurrent entry for %s:\n---\n%s\n---\n", dateStr, body)

	// Generate corrected version
	correctedBody := fixParagraphs(body)
	fmt.Printf("\nCorrected entry (newlines within paragraphs removed):\n---\n%s\n---\n", correctedBody)

	// Prompt for confirmation
	fmt.Print("Update entry with corrected version? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if response != "y" {
		fmt.Println("Update cancelled.")
		return nil
	}

	// Update the entry
	_, err = db.Exec(`
        UPDATE entries_entry 
        SET body = ? 
        WHERE user_id = ? AND "when" = ?`, correctedBody, userID, dateStr)
	if err != nil {
		logger.Printf("Failed to update entry for user %d on %s: %v", userID, dateStr, err)
		return err
	}
	logger.Printf("Updated entry for user %d on %s", userID, dateStr)

	// Requery to confirm update
	var updatedBody string
	err = db.QueryRow(`
        SELECT body 
        FROM entries_entry 
        WHERE user_id = ? AND "when" = ?`, userID, dateStr).Scan(&updatedBody)
	if err != nil {
		logger.Printf("Failed to requery entry for user %d on %s: %v", userID, dateStr, err)
		return err
	}

	fmt.Printf("\nConfirmed updated entry for %s:\n---\n%s\n---\n", dateStr, updatedBody)
	return nil
}

// listEntryDatesF queries and prints all entry dates for user_id=1 in YYYY-MM-DD format.
func listEntryDatesF(db *sql.DB, logger *log.Logger) error {
	const userID = 1
	rows, err := db.Query(`
        SELECT strftime('%Y-%m-%d', "when") 
        FROM entries_entry 
        WHERE user_id = ? 
        ORDER BY "when"`, userID)
	if err != nil {
		logger.Printf("Failed to query entry dates for user %d: %v", userID, err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dateStr string
		if err := rows.Scan(&dateStr); err != nil {
			logger.Printf("Failed to scan date: %v", err)
			return err
		}
		fmt.Println(dateStr)
	}
	if err := rows.Err(); err != nil {
		logger.Printf("Error iterating rows: %v", err)
		return err
	}
	return nil
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	err := godotenv.Load()
	if err != nil {
		logger.Printf("Could not load .env file: %v (continuing without it)", err)
	}

	config := getConfig()

	if config.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDSN,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		defer sentry.Flush(2 * time.Second)
		logger.Println("Sentry is enabled.")
	} else {
		logger.Println("Sentry is disabled.")
	}

	if config.RequiredToAddress == "" {
		log.Fatal("REQUIRED_TO_ADDRESS not set.")
	}

	if config.ReplyToAddress == "" {
		log.Fatal("REPLY_TO_ADDRESS not set.")
	}

	if config.MattEmailAddress == "" {
		log.Fatal("MATT_EMAIL_ADDRESS not set.")
	}

	if config.JournalUser == "" {
		log.Fatal("JOURNAL_USER not set.")
	}
	if config.JournalPassword == "" {
		log.Fatal("JOURNAL_PASSWORD not set.")
	}

	dbPath := "./db.sqlite3"
	if config.DatabaseDirectory != "" {
		dbPath = filepath.Join(config.DatabaseDirectory, "db.sqlite3")
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	logger.Printf("Opened database at %s.", dbPath)

	if config.AWSAccessKeyID == "" {
		log.Fatal("AWS_ACCESS_KEY_ID not set.")
	}
	if config.AWSSecretAccessKey == "" {
		log.Fatal("AWS_SECRET_ACCESS_KEY not set.")
	}
	emailGateway := entries.NewAmazonSESGateway(config)

	// Parse command-line flags
	emailDate := flag.String("email", "", "Send an email prompt for the specified date (YYYY-MM-DD)")
	fixEntryDate := flag.String("fix-entry-date", "", "Fix paragraph newlines for the entry on the specified date (YYYY-MM-DD)")
	listEntryDates := flag.Bool("list-entry-dates", false, "List all entry dates for user_id=1 in YYYY-MM-DD format")
	flag.Parse()

	if *listEntryDates {
		logger.Println("Listing all entry dates")
		err := listEntryDatesF(db, logger)
		if err != nil {
			logger.Fatalf("Failed to list entry dates: %v", err)
		}
		logger.Println("Entry dates listed successfully")
		os.Exit(0)
	}

	if *fixEntryDate != "" {
		// Parse the provided date
		date, err := time.Parse("2006-01-02", *fixEntryDate)
		if err != nil {
			logger.Fatalf("Invalid date format for --fix-entry-date: %v (use YYYY-MM-DD)", err)
		}
		// Process the entry for the specified date
		logger.Printf("Processing entry for %s", *fixEntryDate)
		err = fixEntryForDate(db, date, logger)
		if err != nil {
			logger.Fatalf("Failed to process entry: %v", err)
		}
		logger.Println("Entry processing completed")
		os.Exit(0)
	}

	if *emailDate != "" {
		// Parse the provided date
		date, err := time.Parse("2006-01-02", *emailDate)
		if err != nil {
			logger.Fatalf("Invalid date format for --email: %v (use YYYY-MM-DD)", err)
		}
		// Send email for the specified date
		logger.Printf("Sending email prompt for %s", *emailDate)
		err = entries.SendEmailForDate(db, emailGateway, config, logger, date)
		if err != nil {
			logger.Fatalf("Failed to send email prompt: %v", err)
		}
		logger.Println("Email prompt sent successfully")
		os.Exit(0)
	}

	// Normal server operation if no flags
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/up", up)
	mux.HandleFunc("/journal", journalHandler(db, config, logger))
	username, password := getWebhookAuth(config)
	processor := entries.MakeEmailContentProcessor(config, db, logger)
	mux.HandleFunc("/ses-webhook", webhook.SESWebhookHandler(username, password, processor, logger))

	entries.RunDailyEmailTask(db, emailGateway, config, logger)

	// Start the server
	logger.Println("Server starting on port 8000...")
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		logger.Printf("Server failed to start: %v\n", err)
	}
}
