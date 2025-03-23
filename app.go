package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed go_templates
var templates embed.FS
var tmpl = template.Must(template.ParseFS(templates, "go_templates/index.html"))

func index(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func up(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/up", up)

	fmt.Println("Server starting on port 8000...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
