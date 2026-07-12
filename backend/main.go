package main

import (
	"embed"
	"log"

	"ainexusrouter-core/server"
)

//go:embed public/*
var content embed.FS

func main() {
	port := "20128"
	
	if err := server.RunServer(port, content); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
