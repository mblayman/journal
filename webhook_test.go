package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHappyPath(t *testing.T) {
	username := "user"
	password := "123456"
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte("data")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetBasicAuth(username, password)
	recorder := httptest.NewRecorder()
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "test: ", log.LstdFlags)
	handler := webhookHandler(username, password, logger)

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			recorder.Code, http.StatusOK)
	}
}
