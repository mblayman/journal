package model

import (
	"testing"
)

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

On Wed, Mar 26, 2025, 9:00 AM JourneyInbox Journal <journal@mail.journeyinbox.com> wrote:
> Reply to this prompt to update your journal.
>`,
			},
			wantText: `I got up this morning at 8:30 and brushed my teeth, then left to go to Cafe Ibiza to meet with Jared. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`,
		},
	}

	conf := Config{
		RequiredToAddress: "journal@mail.journeyinbox.com",
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.emailContent.Reply(conf)
			if got != tt.wantText {
				t.Errorf("Reply() = %q, want %q", got, tt.wantText)
			}
		})
	}
}
