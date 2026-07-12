package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"ainexusrouter-core/api"
	"ainexusrouter-core/config"
	"ainexusrouter-core/db"
	"ainexusrouter-core/discovery"
)

// RunServer initializes the application context and starts the HTTP API proxy.
// It runs blocking, so it should typically be called in a goroutine when embedded in Wails.
func RunServer(port string, content embed.FS) error {
	// Providers are compiled into the binary (backend/config/providers.json)
	if err := config.LoadProviders(); err != nil {
		log.Printf("Warning: Could not load embedded providers: %v", err)
	} else {
		log.Printf("Loaded embedded providers config (%d providers)", len(config.GlobalConfig.Providers))
	}

	// Initialize SQLite Database
	if err := db.InitDB(); err != nil {
		log.Printf("Warning: Failed to initialize SQLite database: %v", err)
	}

	// Boot up Dynamic Model Discovery
	discovery.RunDiscovery()

	// Initialize the API server with embedded frontend (if any)
	mux := api.NewServer(content)

	// Same as start.bat: free the port if a previous instance is still bound.
	freePort(port)

	fmt.Printf("ClickToAutomate AI Nexus Router (Go API) starting on port %s\n", port)
	return http.ListenAndServe(":"+port, mux)
}
