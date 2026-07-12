package handlers

import (
	"encoding/json"
	"net/http"

	"ainexusrouter-core/db"
)

// HandleCache manages the semantic cache UI
func HandleCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		cacheEntries := db.GetAllCache()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"cache": cacheEntries,
		})
		return
	}

	if r.Method == http.MethodDelete {
		if err := db.ClearCache(); err != nil {
			http.Error(w, "Failed to clear cache", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
