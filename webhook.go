package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type EmailContent struct {
	To      string
	Text    string
	Subject string
}

// EmailContentProcessor is the callback that does all the necessary
// processing on the (relatively) raw email data.
type EmailContentProcessor func(EmailContent)

func webhookHandler(username, password string, processor EmailContentProcessor, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Read the raw request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Save the request body to a file
		filename := fmt.Sprintf("webhook_%s.bin", time.Now().Format("20060102_150405.999"))
		filePath := filepath.Join("/var/db", filename)
		if err := os.WriteFile(filePath, bodyBytes, 0644); err != nil {
			logger.Printf("Error writing request body to file %s: %v", filePath, err)
			// Don’t fail the request, just log the error
		} else {
			logger.Printf("Saved request body to %s", filePath)
		}

		// Reconstruct the body for processing
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "multipart/form-data" {
			logger.Printf("Invalid Content-Type: %v", err)
			http.Error(w, "Expected multipart/form-data", http.StatusBadRequest)
			return
		}

		reader := multipart.NewReader(r.Body, params["boundary"])
		form, err := reader.ReadForm(10 << 20) // 10MB max memory
		if err != nil {
			logger.Printf("Error reading form: %v", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		if emailValues, ok := form.Value["email"]; !ok || len(emailValues) == 0 {
			logger.Printf("No email field found in webhook")
			http.Error(w, "Missing email field", http.StatusBadRequest)
			return
		} else {
			content, err := extractContent(emailValues[0], logger)
			if err != nil {
				logger.Printf("Error extracting content: %v", err)
				http.Error(w, "Error processing email", http.StatusInternalServerError)
				return
			}
			processor(content)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// extractContent pulls the details from the raw email message.
func extractContent(emailRaw string, logger *log.Logger) (EmailContent, error) {
	msg, err := mail.ReadMessage(bytes.NewReader([]byte(emailRaw)))
	if err != nil {
		return EmailContent{}, fmt.Errorf("failed to parse email: %v", err)
	}

	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return EmailContent{}, fmt.Errorf("invalid Content-Type: %v", err)
	}

	if mediaType == "multipart/alternative" {
		boundary := params["boundary"]
		if boundary == "" {
			return EmailContent{}, fmt.Errorf("missing boundary in multipart/alternative")
		}

		body, err := io.ReadAll(msg.Body)
		if err != nil {
			return EmailContent{}, fmt.Errorf("failed to read email body: %v", err)
		}

		var textContent string
		reader := multipart.NewReader(bytes.NewReader(body), boundary)
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return EmailContent{}, fmt.Errorf("error reading multipart part: %v", err)
			}
			partContentType := part.Header.Get("Content-Type")
			logger.Printf("Found part with Content-Type: %s", partContentType)
			if strings.HasPrefix(partContentType, "text/plain") {
				content, err := io.ReadAll(part)
				if err != nil {
					logger.Printf("Error reading text/plain part: %v", err)
					return EmailContent{}, fmt.Errorf("error reading text/plain part: %v", err)
				}
				textContent = string(content)
			}
		}
		if textContent != "" {
			logger.Printf("Extracted text content length: %d", len(textContent))
			return EmailContent{To: to, Subject: subject, Text: textContent}, nil
		}
		return EmailContent{}, fmt.Errorf("no text/plain part found")
	}

	if mediaType == "text/plain" {
		content, err := io.ReadAll(msg.Body)
		if err != nil {
			return EmailContent{}, fmt.Errorf("failed to read text body: %v", err)
		}
		return EmailContent{To: to, Subject: subject, Text: string(content)}, nil
	}

	return EmailContent{}, fmt.Errorf("unsupported Content-Type: %s", mediaType)
}
