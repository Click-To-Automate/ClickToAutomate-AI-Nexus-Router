package main

import (
	"embed"
	"log"

	"ainexusrouter-core/apiport"
	"ainexusrouter-core/server"
)

//go:embed public/*
var content embed.FS

func main() {
	if err := server.RunServer(apiport.Port, "", content); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
