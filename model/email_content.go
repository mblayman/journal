package model

import (
	"strings"
)

// EmailContent is the data extracted from webhook responses.
type EmailContent struct {
	To      string
	Text    string
	Subject string
}

// ToAddress extracts the address portion of the To field
func (emailContent *EmailContent) ToAddress() string {
	start := strings.Index(emailContent.To, "<")
	end := strings.Index(emailContent.To, ">")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return emailContent.To[start+1 : end]
}

// Reply extracts the reply content from the email text, excluding the quoted original message.
// It joins lines within paragraphs and stops before the paragraph containing ToAddress.
func (emailContent *EmailContent) Reply() string {
	lines := strings.Split(emailContent.Text, "\n")
	toAddress := emailContent.ToAddress()
	var paragraphs []string
	var currentParagraph []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			// Blank line ends a paragraph
			if len(currentParagraph) > 0 {
				joinedParagraph := strings.Join(currentParagraph, " ")
				// Check if this paragraph contains toAddress
				if strings.Contains(joinedParagraph, toAddress) {
					// Stop here, excluding this paragraph
					break
				}
				paragraphs = append(paragraphs, joinedParagraph)
				currentParagraph = nil
			}
		} else {
			currentParagraph = append(currentParagraph, trimmedLine)
		}
	}

	return strings.Join(paragraphs, "\n\n")
}

// EmailContentProcessor is the callback that does all the necessary
// processing on the (relatively) raw email data.
type EmailContentProcessor func(EmailContent)
