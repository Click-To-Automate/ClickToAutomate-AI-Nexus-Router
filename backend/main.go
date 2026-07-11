package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed public/*
var content embed.FS

func main() {
	port := "20128"

	// Create a sub-FS for the public directory
	publicFS, err := fs.Sub(content, "public")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	// Setup static file server for frontend
	fileServer := http.FileServer(http.FS(publicFS))
	http.Handle("/", fileServer)

	// Stub API endpoint
	http.HandleFunc("/v1/models", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data": []}`)
	})

	fmt.Printf("ClickToAutomate AI Nexus Router (Go Edition) backend starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
