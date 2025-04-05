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
