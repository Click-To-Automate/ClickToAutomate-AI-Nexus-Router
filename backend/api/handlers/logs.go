package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// LogEntry represents a single log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}

// HandleLogs returns mock log entries for demonstration.
// Replace this with actual log fetching logic (e.g., from a file or database).
func HandleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock log entries for demonstration
	logs := []LogEntry{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:    "INFO",
			Message:  "Server started on port 8080",
		},
		{
			Timestamp: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			Level:    "WARN",
			Message:  "High memory usage detected",
		},
		{
			Timestamp: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			Level:    "ERROR",
			Message:  "Failed to connect to provider: OpenAI",
		},
	}

	json.NewEncoder(w).Encode(logs)
}