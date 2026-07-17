package handlers

import (
	"encoding/json"
	"net/http"

	"ainexusrouter-core/db"
)

// HandleUsage returns the current usage statistics for all providers
func HandleUsage(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers for the frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	usage := db.GetAllUsage()
	latencies := db.GetLatencies()
	
	type UsageResponse struct {
		Usage     map[string]db.UsageData `json:"usage"`
		Latencies map[string]int          `json:"latencies"`
	}
	
	resp := UsageResponse{
		Usage:     usage,
		Latencies: latencies,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode usage", http.StatusInternalServerError)
		return
	}
}
