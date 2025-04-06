package model

import (
	"testing"
)

// TestToAddress checks the address portion of the To field.
func TestToAddress(t *testing.T) {
	emailContent := EmailContent{
		To: "JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>",
	}
	expectedAddress := "journal.abcdef1@email.journeyinbox.com"

	if emailContent.ToAddress() != expectedAddress {
		t.Errorf("Expected %q, got %q", expectedAddress, emailContent.ToAddress())
	}
}

// TestReply tests the extraction of reply text before the quoted section using ToAddress.
func TestReply(t *testing.T) {
	tests := []struct {
		name         string
		emailContent EmailContent
		wantText     string
	}{
		{
			name: "Valid text with quoted section",
			emailContent: EmailContent{
				To: "JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>",
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
			wantText: `I got up this morning at 8:30 and brushed my teeth, then left to go to Cafe Ibiza to meet with Jared. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
		},
		{
			name: "Text without quoted section",
			emailContent: EmailContent{
				To:   "JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>",
				Text: "Just some text without a quote.",
			},
			wantText: "Just some text without a quote.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.emailContent.Reply()
			if got != tt.wantText {
				t.Errorf("Reply() = %q, want %q", got, tt.wantText)
			}
		})
	}
}
