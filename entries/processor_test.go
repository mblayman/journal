package entries

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/mblayman/journal/model"
)

// TestMakeEmailContentProcessorToAddress tests the To address validation in the processor with a valid Subject.
func TestMakeEmailContentProcessorToAddress(t *testing.T) {
	tests := []struct {
		name         string
		requiredTo   string
		emailContent model.EmailContent
		expectedLog  string // Expected log output (empty for invalid To, specific for valid To)
	}{
		{
			name:       "Matching To address with valid Subject",
			requiredTo: "journal.abcdef1@email.journeyinbox.com",
			emailContent: model.EmailContent{
				To:      "JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>",
				Subject: "It's Wednesday, Mar. 26, 2025. How are you?",
			},
			expectedLog: "Parsed date from Subject: 2025-03-26T00:00:00Z",
		},
		{
			name:       "Non-matching To address with valid Subject",
			requiredTo: "journal.abcdef1@email.journeyinbox.com",
			emailContent: model.EmailContent{
				To:      "JourneyInbox Journal <journal.xyz789@email.journeyinbox.com>",
				Subject: "It's Wednesday, Mar. 26, 2025. How are you?",
			},
			expectedLog: "Invalid To address: journal.xyz789@email.journeyinbox.com",
		},
		{
			name:       "Invalid To format with valid Subject",
			requiredTo: "journal.abcdef1@email.journeyinbox.com",
			emailContent: model.EmailContent{
				To:      "Invalid Format",
				Subject: "It's Wednesday, Mar. 26, 2025. How are you?",
			},
			expectedLog: "Invalid To address:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", log.LstdFlags)

			processor := MakeEmailContentProcessor(tt.requiredTo, logger)
			processor(tt.emailContent)

			logOutput := strings.TrimSpace(logBuf.String())
			if !strings.Contains(logOutput, tt.expectedLog) {
				t.Errorf("Expected log output to contain %q, got %q", tt.expectedLog, logOutput)
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
