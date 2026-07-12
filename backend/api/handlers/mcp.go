package handlers

import (
	"encoding/json"
	"net/http"
)

// HandleMCP manages the Model Context Protocol local agents whitelist
func HandleMCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		// Mock MCP whitelist for now
		json.NewEncoder(w).Encode(map[string]interface{}{
			"allowed_commands": []string{"ls", "git status"},
			"allowed_dirs":     []string{"/tmp", "./src"},
		})
		return
	}

	if r.Method == http.MethodPost {
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "MCP settings saved (stub)"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
