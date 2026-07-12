package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"ainexusrouter-core/apiport"
	"ainexusrouter-core/server"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	
	// Start the backend API proxy in a goroutine (dev: 20128, prod: via ldflags)
	// while Wails handles the GUI natively
	go func() {
		// Provide an empty embed.FS to the server since Wails handles assets now
		var emptyEmbed embed.FS
		if err := server.RunServer(apiport.Port, emptyEmbed); err != nil {
			log.Fatalf("API Server crashed: %v", err)
		}
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ClickToAutomate AI Nexus Router",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 248, G: 250, B: 252, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
