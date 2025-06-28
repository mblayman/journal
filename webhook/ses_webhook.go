package webhook

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

// SESWebhookPayload represents the SNS notification structure.
type SESWebhookPayload struct {
	NotificationType string `json:"notificationType"`
	Mail             struct {
		Source      string   `json:"source"`
		Destination []string `json:"destination"`
		Subject     string   `json:"subject"`
	} `json:"mail"`
	Receipt struct {
		Action struct {
			Type     string `json:"type"`
			TopicArn string `json:"topicArn"`
			Encoding string `json:"encoding"`
		} `json:"action"`
	} `json:"receipt"`
	Content string `json:"content"`
}

// SESWebhookHandler processes SNS notifications from AWS SES.
func SESWebhookHandler(username, password string, processor model.EmailContentProcessor, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Printf("Invalid method: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			logger.Printf("Unauthorized access attempt")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Read and parse SNS payload
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request", http.StatusBadRequest)
			return
		}

		// Check for SNS SubscriptionConfirmation
		var snsMessage map[string]interface{}
		if err := json.Unmarshal(body, &snsMessage); err == nil {
			if snsMessage["Type"] == "SubscriptionConfirmation" {
				subscribeURL, ok := snsMessage["SubscribeURL"].(string)
				if ok {
					logger.Printf("Confirming SNS subscription: %s", subscribeURL)
					resp, err := http.Get(subscribeURL)
					if err != nil {
						logger.Printf("Failed to confirm SNS subscription: %v", err)
						http.Error(w, "Error confirming subscription", http.StatusInternalServerError)
						return
					}
					defer resp.Body.Close()
					logger.Printf("SNS subscription confirmed")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ok"))
					return
				}
			}
		}

		var payload SESWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			logger.Printf("Error parsing JSON payload: %v, body: %s", err, string(body))
			http.Error(w, "Error parsing JSON\n", http.StatusBadRequest)
			return
		}

		// Verify notification type and action
		if payload.NotificationType != "Received" || payload.Receipt.Action.Type != "SNS" {
			logger.Printf("Invalid notification type: %s or action: %s", payload.NotificationType, payload.Receipt.Action.Type)
			http.Error(w, "Invalid notification", http.StatusBadRequest)
			return
		}

		// Decode base64 content
		emailRaw, err := base64.StdEncoding.DecodeString(payload.Content)
		if err != nil {
			logger.Printf("Error decoding base64 content: %v", err)
			http.Error(w, "Error decoding email content", http.StatusBadRequest)
			return
		}

		// Parse email content
		content, err := extractSESContent(emailRaw, payload, logger)
		if err != nil {
			logger.Printf("Error extracting content: %v", err)
			http.Error(w, "Error processing email", http.StatusInternalServerError)
			return
		}

		processor(content)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// extractSESContent extracts To, Subject, and text body from the raw email, using payload.Destination as fallback.
func extractSESContent(emailRaw []byte, payload SESWebhookPayload, logger *log.Logger) (model.EmailContent, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(emailRaw))
	if err != nil {
		return model.EmailContent{}, fmt.Errorf("failed to parse email: %v", err)
	}

	// Extract To header without modifying case
	to := msg.Header.Get("To")
	if to == "" && len(payload.Mail.Destination) > 0 {
		to = payload.Mail.Destination[0]
		logger.Printf("Using fallback To address from payload: %q", to)
	}
	if to == "" {
		logger.Printf("No valid To address found in header: %q or payload: %v", msg.Header.Get("To"), payload.Mail.Destination)
		return model.EmailContent{}, fmt.Errorf("no valid To address found")
	}

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
				// Handle quoted-printable encoding
				if part.Header.Get("Content-Transfer-Encoding") == "quoted-printable" {
					decoded, err := decodeQuotedPrintable(string(content))
					if err != nil {
						logger.Printf("Error decoding quoted-printable: %v", err)
						return model.EmailContent{}, fmt.Errorf("error decoding quoted-printable: %v", err)
					}
					textContent = decoded
				} else {
					textContent = string(content)
				}
			}
		}
		if textContent != "" {
			return model.EmailContent{To: to, Subject: subject, Text: textContent}, nil
		}
		return model.EmailContent{}, fmt.Errorf("no text/plain part found")
	}

	return model.EmailContent{}, fmt.Errorf("unsupported Content-Type: %s", mediaType)
}

// decodeQuotedPrintable decodes quoted-printable text (simplified for this use case).
func decodeQuotedPrintable(input string) (string, error) {
	// Replace =XX with the corresponding byte
	var result strings.Builder
	i := 0
	for i < len(input) {
		if i+2 < len(input) && input[i] == '=' && isHexChar(input[i+1]) && isHexChar(input[i+2]) {
			hex := input[i+1 : i+3]
			b, err := hexToByte(hex)
			if err != nil {
				return "", err
			}
			result.WriteByte(b)
			i += 3
		} else {
			result.WriteByte(input[i])
			i++
		}
	}
	return result.String(), nil
}

// isHexChar checks if a character is a valid hexadecimal digit.
func isHexChar(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')
}

// hexToByte converts a two-character hex string to a byte.
func hexToByte(hex string) (byte, error) {
	var b byte
	_, err := fmt.Sscanf(hex, "%x", &b)
	return b, err
}
