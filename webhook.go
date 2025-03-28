package main

import (
	"io"
	"log"
	"net/http"
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Printf("Error reading request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		logger.Printf("Request body: %s", string(body))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
