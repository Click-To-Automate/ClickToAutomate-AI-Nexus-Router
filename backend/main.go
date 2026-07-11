package main

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

//go:embed public/*
var content embed.FS

func main() {
	port := "20128"

	// Load dynamic providers configuration
	if err := config.LoadProviders("config/providers.json"); err != nil {
		log.Printf("Warning: Could not load providers.json, using defaults: %v", err)
	}

	// Initialize SQLite Database
	if err := db.InitDB(); err != nil {
		log.Printf("Warning: Failed to initialize SQLite database: %v", err)
	}

	// Boot up Dynamic Model Discovery
	discovery.RunDiscovery()

	// Initialize the API server with embedded frontend
	mux := api.NewServer(content)

	fmt.Printf("ClickToAutomate AI Nexus Router (Go Edition) backend starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
