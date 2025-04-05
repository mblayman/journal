package model

import "strings"

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
