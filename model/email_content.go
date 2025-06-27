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

// Reply extracts the reply content from the email text, excluding the quoted original message.
func (emailContent *EmailContent) Reply(conf Config) string {
	lines := strings.Split(emailContent.Text, "\n")
	toAddress := conf.RequiredToAddress
	var paragraphs []string
	var currentParagraph []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			// Blank line ends a paragraph
			if len(currentParagraph) > 0 {
				joinedParagraph := strings.Join(currentParagraph, " ")
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
