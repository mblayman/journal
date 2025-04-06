package webhook

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"strings"

	"github.com/mblayman/journal/model"
)

// WebhookHandler returns a handler func that can process webhook
// data sent by SendGrid. The email content is delegated to the processor.
func WebhookHandler(username, password string, processor model.EmailContentProcessor, logger *log.Logger) http.HandlerFunc {
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
func extractContent(emailRaw string, logger *log.Logger) (model.EmailContent, error) {
	msg, err := mail.ReadMessage(bytes.NewReader([]byte(emailRaw)))
	if err != nil {
		return model.EmailContent{}, fmt.Errorf("failed to parse email: %v", err)
	}

	to := msg.Header.Get("To")
	subject := msg.Header.Get("Subject")

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return model.EmailContent{}, fmt.Errorf("invalid Content-Type: %v", err)
	}

	if mediaType == "multipart/alternative" {
		boundary := params["boundary"]
		if boundary == "" {
			return model.EmailContent{}, fmt.Errorf("missing boundary in multipart/alternative")
		}

		body, err := io.ReadAll(msg.Body)
		if err != nil {
			return model.EmailContent{}, fmt.Errorf("failed to read email body: %v", err)
		}

		var textContent string
		reader := multipart.NewReader(bytes.NewReader(body), boundary)
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return model.EmailContent{}, fmt.Errorf("error reading multipart part: %v", err)
			}
			partContentType := part.Header.Get("Content-Type")
			if strings.HasPrefix(partContentType, "text/plain") {
				content, err := io.ReadAll(part)
				if err != nil {
					logger.Printf("Error reading text/plain part: %v", err)
					return model.EmailContent{}, fmt.Errorf("error reading text/plain part: %v", err)
				}
				textContent = string(content)
			}
		}
		if textContent != "" {
			return model.EmailContent{To: to, Subject: subject, Text: textContent}, nil
		}
		return model.EmailContent{}, fmt.Errorf("no text/plain part found")
	}

	return model.EmailContent{}, fmt.Errorf("unsupported Content-Type: %s", mediaType)
}
