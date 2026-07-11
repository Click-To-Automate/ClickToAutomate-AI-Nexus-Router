package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"ainexusrouter-core/api/handlers"
)

// NewServer initializes the HTTP multiplexer with all routes
func NewServer(frontendFS embed.FS) *http.ServeMux {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/v1/chat/completions", handlers.HandleChatCompletions)
	mux.HandleFunc("/v1/models", handlers.HandleModels)
	mux.HandleFunc("/v1/keys", handlers.HandleKeys)
	mux.HandleFunc("/v1/usage", handlers.HandleUsage)

	// Frontend Static Files
	publicFS, err := fs.Sub(frontendFS, "public")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}
	fileServer := http.FileServer(http.FS(publicFS))

	// We route root to the static file server
	mux.Handle("/", fileServer)

	return mux
}
