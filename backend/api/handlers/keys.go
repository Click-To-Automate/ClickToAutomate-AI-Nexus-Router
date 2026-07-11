package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"ainexusrouter-core/db"
)

type SetKeyRequest struct {
	ProviderID string `json:"provider_id"`
	APIKey     string `json:"api_key"`
}

// HandleKeys provides GET (list keys) and POST (upsert key)
func HandleKeys(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers for the frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		// For Phase 3, we simply return a 200 OK array.
		// A full list endpoint would iterate through the DB,
		// but since modernc.org/sqlite doesn't have an easy ORM loaded yet,
		// we can query all rows.
		
		if db.DB == nil {
			http.Error(w, `{"error": "DB not initialized"}`, http.StatusInternalServerError)
			return
		}

		rows, err := db.DB.Query("SELECT provider_id, api_key FROM provider_keys")
		if err != nil {
			http.Error(w, `{"error": "Failed to query keys"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]string
		for rows.Next() {
			var pid, key string
			if err := rows.Scan(&pid, &key); err == nil {
				// Mask the key before sending to frontend
				masked := ""
				if len(key) > 8 {
					masked = key[:4] + "..." + key[len(key)-4:]
				} else {
					masked = "***"
				}
				results = append(results, map[string]string{
					"provider_id": pid,
					"masked_key":  masked,
					"is_set":      "true",
				})
			}
		}
		
		if results == nil {
			results = []map[string]string{}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"keys": results,
		})
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req SetKeyRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.ProviderID == "" || req.APIKey == "" {
			http.Error(w, "provider_id and api_key are required", http.StatusBadRequest)
			return
		}

		if err := db.SetKey(req.ProviderID, req.APIKey); err != nil {
			http.Error(w, "Failed to save key", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
