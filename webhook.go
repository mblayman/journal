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
	"strings"
)

func webhookHandler(username, password string, logger *log.Logger) http.HandlerFunc {
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
		defer r.Body.Close()

		var to string
		if toValues, ok := form.Value["to"]; !ok || len(toValues) == 0 {
			logger.Printf("No to field found in webhook")
			http.Error(w, "Missing to field", http.StatusBadRequest)
			return
		} else {
			to = toValues[0]
		}

		if emailValues, ok := form.Value["email"]; !ok || len(emailValues) == 0 {
			logger.Printf("No email field found in webhook")
			http.Error(w, "Missing email field", http.StatusBadRequest)
			return
		} else {
			text, err := extractTextContent(emailValues[0], logger)
			if err != nil {
				logger.Printf("Error extracting text: %v", err)
				http.Error(w, "Error processing email", http.StatusInternalServerError)
				return
			}
			logger.Printf("To: %s", to)
			// Replace newlines with spaces to keep it one line
			singleLineText := strings.ReplaceAll(text, "\n", " ")
			logger.Printf("Text Content: %s", singleLineText)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func extractTextContent(emailRaw string, logger *log.Logger) (string, error) {
	msg, err := mail.ReadMessage(bytes.NewReader([]byte(emailRaw)))
	if err != nil {
		return "", fmt.Errorf("failed to parse email: %v", err)
	}

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("invalid Content-Type: %v", err)
	}

	if mediaType == "multipart/alternative" {
		boundary := params["boundary"]
		if boundary == "" {
			return "", fmt.Errorf("missing boundary in multipart/alternative")
		}

		body, err := io.ReadAll(msg.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read email body: %v", err)
		}

		var textContent string
		reader := multipart.NewReader(bytes.NewReader(body), boundary)
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", fmt.Errorf("error reading multipart part: %v", err)
			}
			partContentType := part.Header.Get("Content-Type")
			logger.Printf("Found part with Content-Type: %s", partContentType)
			if strings.HasPrefix(partContentType, "text/plain") {
				content, err := io.ReadAll(part)
				if err != nil {
					return "", fmt.Errorf("error reading text/plain part: %v", err)
				}
				textContent = string(content)
			}
		}
		if textContent != "" {
			return textContent, nil
		}
		return "", fmt.Errorf("no text/plain part found")
	}

	return "", fmt.Errorf("unsupported Content-Type: %s", mediaType)
}
