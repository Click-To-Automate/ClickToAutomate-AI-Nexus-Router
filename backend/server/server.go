package server

import (
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ainexusrouter-core/api"
	"ainexusrouter-core/api/handlers"
	"ainexusrouter-core/config"
	"ainexusrouter-core/db"
	"ainexusrouter-core/discovery"
)

// RunServer initializes the application context and starts the HTTP API proxy.
// It runs blocking, so it should typically be called in a goroutine when embedded in Wails.
func RunServer(port string, dbPath string, content embed.FS) error {
	// Initialize SQLite Database first so we get the default dbPath if it was empty
	if err := db.InitDB(dbPath); err != nil {
		log.Printf("Warning: Failed to initialize SQLite database: %v", err)
	}

	// Create a log file in the same directory as the database
	actualDBPath := db.DBPath
	logFileDir := filepath.Dir(actualDBPath)
	logFilePath := filepath.Join(logFileDir, "logs.txt")
	
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Warning: Could not open log file %s: %v", logFilePath, err)
		// Fallback without file logging
		multiWriter := io.MultiWriter(os.Stdout, handlers.GlobalLogBuffer)
		log.SetOutput(multiWriter)
	} else {
		// Hook up the global logger with file
		multiWriter := io.MultiWriter(os.Stdout, handlers.GlobalLogBuffer, logFile)
		log.SetOutput(multiWriter)
	}

	// Providers are compiled into the binary (backend/config/providers.json)
	if err := config.LoadProviders(); err != nil {
		log.Printf("Warning: Could not load embedded providers: %v", err)
	} else {
		log.Printf("Loaded embedded providers config (%d providers)", len(config.GlobalConfig.Providers))
	}



	// Boot up Dynamic Model Discovery
	discovery.RunDiscovery()

	// Initialize the API server with embedded frontend (if any)
	mux := api.NewServer(content)

	// Spawn background cache sweeper (TTL: 30 days)
	go func() {
		// Run once on startup
		if err := db.CleanExpiredCache(30); err != nil {
			log.Printf("Warning: Failed to clean expired cache on startup: %v", err)
		}
		
		// Then run every 24 hours
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := db.CleanExpiredCache(30); err != nil {
				log.Printf("Warning: Failed to clean expired cache: %v", err)
			} else {
				log.Printf("Successfully ran daily cache cleanup")
			}
		}
	}()

	// Same as start.bat: free the port if a previous instance is still bound.
	freePort(port)

	log.Printf("ClickToAutomate AI Nexus Router (Go API) starting on port %s\n", port)
	return http.ListenAndServe(":"+port, mux)
}
