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
mux.HandleFunc("/v1/providers", handlers.HandleProviders)
mux.HandleFunc("/v1/settings", handlers.HandleSettings)
mux.HandleFunc("/v1/cache", handlers.HandleCache)
mux.HandleFunc("/v1/mcp", handlers.HandleMCP)
mux.HandleFunc("/v1/logs", handlers.HandleLogs)

	// Frontend Static Files (optional when Wails serves the UI)
	publicFS, err := fs.Sub(frontendFS, "public")
	if err != nil {
		log.Printf("No embedded public assets (Wails mode): %v", err)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok","service":"AI Nexus Router API"}`))
		})
		return mux
	}
	fileServer := http.FileServer(http.FS(publicFS))
	mux.Handle("/", fileServer)

	return mux
}
