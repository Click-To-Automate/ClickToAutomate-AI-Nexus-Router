package handlers

import (
	"encoding/json"
	"net/http"

	"ainexusrouter-core/apiport"
	"ainexusrouter-core/db"
)

// HandleSettings manages app settings like compression toggles, dark mode, etc.
func HandleSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		settings := map[string]string{
			"rtk_engine":     db.GetSetting("rtk_engine", "true"),
			"caveman_engine": db.GetSetting("caveman_engine", "false"),
			"semantic_cache": db.GetSetting("semantic_cache", "false"),
			"port":           db.GetSetting("port", apiport.Port),
			"dark_mode":      db.GetSetting("dark_mode", "true"),
		}
		json.NewEncoder(w).Encode(settings)
		return
	}

	if r.Method == http.MethodPost {
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for key, value := range payload {
			db.SetSetting(key, value)
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
