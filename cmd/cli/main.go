package main

import (
	"embed"
	"flag"
	"log"

	"ainexusrouter-core/apiport"
	"ainexusrouter-core/server"
)

func main() {
	var port string
	var dbPath string

	flag.StringVar(&port, "port", apiport.Port, "Port to run the API server on")
	flag.StringVar(&dbPath, "db", "", "Path to the sqlite database file")
	flag.Parse()

	log.Printf("Starting CLI with port=%s, db=%s", port, dbPath)

	var emptyEmbed embed.FS

	if err := server.RunServer(port, dbPath, emptyEmbed); err != nil {
		log.Fatalf("API Server crashed: %v", err)
	}

	// Wait forever
	select {}
}
