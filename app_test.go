package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mblayman/journal/model"

	_ "modernc.org/sqlite"
)

func TestEditHandlerGET(t *testing.T) {
	// Create in-memory DB
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Create table and insert test data
	_, err = db.Exec(`CREATE TABLE entries_entry ("when" TEXT, body TEXT, user_id INTEGER)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`INSERT INTO entries_entry ("when", body, user_id) VALUES (?, ?, ?)`, "2023-01-01", "Test entry", 1)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Mock config and logger
	config := model.Config{JournalUser: "testuser", JournalPassword: "testpass"}
	logger := log.New(nil, "", 0) // Discard logs

	handler := editHandler(db, config, logger)

	// Test GET request
	req := httptest.NewRequest("GET", "/edit?date=2023-01-01", nil)
	req.SetBasicAuth("testuser", "testpass")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Test entry") {
		t.Errorf("Expected response to contain 'Test entry', got %s", w.Body.String())
	}
}

func TestEditHandlerPOST(t *testing.T) {
	// Create in-memory DB
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Create table and insert test data
	_, err = db.Exec(`CREATE TABLE entries_entry ("when" TEXT, body TEXT, user_id INTEGER)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`INSERT INTO entries_entry ("when", body, user_id) VALUES (?, ?, ?)`, "2023-01-01", "Old entry", 1)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Mock config and logger
	config := model.Config{JournalUser: "testuser", JournalPassword: "testpass"}
	logger := log.New(nil, "", 0)

	handler := editHandler(db, config, logger)

	// Test POST request
	form := strings.NewReader("body=Updated entry")
	req := httptest.NewRequest("POST", "/edit?date=2023-01-01", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("testuser", "testpass")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status 303, got %d", w.Code)
	}

	// Check if DB was updated
	var body string
	err = db.QueryRow(`SELECT body FROM entries_entry WHERE "when" = ? AND user_id = 1`, "2023-01-01").Scan(&body)
	if err != nil {
		t.Fatalf("Failed to query DB: %v", err)
	}
	if body != "Updated entry" {
		t.Errorf("Expected body 'Updated entry', got '%s'", body)
	}
}

func TestEditAndJournalE2E(t *testing.T) {
	// Create in-memory DB
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Create table and insert initial entry
	_, err = db.Exec(`CREATE TABLE entries_entry ("when" TEXT, body TEXT, user_id INTEGER)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`INSERT INTO entries_entry ("when", body, user_id) VALUES (?, ?, ?)`, "2023-01-01", "Initial entry.", 1)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Mock config and logger
	config := model.Config{JournalUser: "testuser", JournalPassword: "testpass"}
	logger := log.New(nil, "", 0)

	// Test POST to edit with double newlines (using \r\n like browser)
	editHandler := editHandler(db, config, logger)
	form := strings.NewReader("body=First paragraph.\r\n\r\nSecond paragraph.")
	req := httptest.NewRequest("POST", "/edit?date=2023-01-01", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("testuser", "testpass")
	w := httptest.NewRecorder()
	editHandler.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect 303, got %d", w.Code)
	}

	// Check DB was updated (normalized)
	var body string
	err = db.QueryRow(`SELECT body FROM entries_entry WHERE "when" = ? AND user_id = 1`, "2023-01-01").Scan(&body)
	if err != nil {
		t.Fatalf("Failed to query DB: %v", err)
	}
	if body != "First paragraph.\n\nSecond paragraph." {
		t.Errorf("Expected normalized body with double newline, got %q", body)
	}

	// Now test GET to journal
	journalHandler := journalHandler(db, config, logger)
	req2 := httptest.NewRequest("GET", "/journal", nil)
	req2.SetBasicAuth("testuser", "testpass")
	w2 := httptest.NewRecorder()
	journalHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200 for journal, got %d", w.Code)
	}

	html := w2.Body.String()
	// Check for two <p> tags in the HTML
	pCount := strings.Count(html, "<p>")
	if pCount < 2 {
		t.Errorf("Expected at least 2 <p> tags for multiple paragraphs, got %d. HTML: %s", pCount, html)
	}

	// Check that <br> is not present (since double newlines should split)
	if strings.Contains(html, "<br>") {
		t.Errorf("Unexpected <br> in HTML, should be separate <p>. HTML: %s", html)
	}
}

func TestEditAndJournalE2ESingleNewline(t *testing.T) {
	// Create in-memory DB
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`CREATE TABLE entries_entry ("when" TEXT, body TEXT, user_id INTEGER)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = db.Exec(`INSERT INTO entries_entry ("when", body, user_id) VALUES (?, ?, ?)`, "2023-01-01", "Initial.", 1)
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	// Mock config and logger
	config := model.Config{JournalUser: "testuser", JournalPassword: "testpass"}
	logger := log.New(nil, "", 0)

	// Test POST to edit with single newline (using \r\n)
	editHandler := editHandler(db, config, logger)
	form := strings.NewReader("body=First line.\r\nSecond line.")
	req := httptest.NewRequest("POST", "/edit?date=2023-01-01", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("testuser", "testpass")
	w := httptest.NewRecorder()
	editHandler.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect 303, got %d", w.Code)
	}

	// Check DB (normalized to single paragraph)
	var body string
	err = db.QueryRow(`SELECT body FROM entries_entry WHERE "when" = ? AND user_id = 1`, "2023-01-01").Scan(&body)
	if err != nil {
		t.Fatalf("Failed to query DB: %v", err)
	}
	if body != "First line. Second line." {
		t.Errorf("Expected normalized body with space, got %q", body)
	}

	// Test GET to journal
	journalHandler := journalHandler(db, config, logger)
	req2 := httptest.NewRequest("GET", "/journal", nil)
	req2.SetBasicAuth("testuser", "testpass")
	w2 := httptest.NewRecorder()
	journalHandler.ServeHTTP(w2, req2)

	html := w2.Body.String()
	// Should have two <p> (date + body), body normalized so no <br>
	pCount := strings.Count(html, "<p>")
	if pCount != 2 {
		t.Errorf("Expected 2 <p> tags (date + body), got %d. HTML: %s", pCount, html)
	}
	if strings.Contains(html, "<br>") {
		t.Errorf("Unexpected <br> in HTML. HTML: %s", html)
	}
	// Check body is normalized
	if !strings.Contains(html, "First line. Second line.") {
		t.Errorf("Expected normalized body in HTML. HTML: %s", html)
	}
}
