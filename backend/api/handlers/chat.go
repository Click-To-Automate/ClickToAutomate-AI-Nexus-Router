package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"ainexusrouter-core/config"
	"ainexusrouter-core/db"
	"ainexusrouter-core/discovery"
	"ainexusrouter-core/proxy"
)

// OpenAI-compatible Chat Completion Request
type ChatCompletionRequest struct {
	Model    string                   `json:"model"`
	Messages []map[string]interface{} `json:"messages"`
	Stream   bool                     `json:"stream,omitempty"`
}

func HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var chatReq ChatCompletionRequest
	if err := json.Unmarshal(body, &chatReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if chatReq.Model == "" {
		http.Error(w, "Model is required", http.StatusBadRequest)
		return
	}

	// Log the incoming call with relevant details (avoiding huge message bodies)
	log.Printf("[Router] Incoming Request -> Model: %s | Messages: %d | Stream: %v\n", chatReq.Model, len(chatReq.Messages), chatReq.Stream)

	// Delegate to the proxy engine
	err = proxy.HandleProxyRequest(w, r, body, chatReq.Model, chatReq.Stream)
	if err != nil {
		fmt.Printf("Proxy Error: %v\n", err)
	}
}

type ProviderGroup struct {
	ProviderID   string   `json:"provider_id"`
	ProviderName string   `json:"provider_name"`
	Models       []string `json:"models"`
}

type ModelsResponse struct {
	Object string          `json:"object"`
	Data   []ProviderGroup `json:"data"`
}

func HandleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	resp := ModelsResponse{
		Object: "list",
		Data:   []ProviderGroup{},
	}

	activeKeys := make(map[string]bool)
	if db.DB != nil {
		rows, err := db.DB.Query("SELECT DISTINCT provider_id FROM provider_keys_multi")
		if err == nil {
			for rows.Next() {
				var pid string
				if err := rows.Scan(&pid); err == nil {
					activeKeys[pid] = true
				}
			}
			rows.Close()
		}
	}

	discovery.Mu.RLock()
	for _, p := range config.GlobalConfig.Providers {
		if !activeKeys[p.ID] {
			continue
		}

		models, ok := discovery.ProviderModels[p.ID]
		if !ok || len(models) == 0 {
			// Fallback if discovery is pending or failed (e.g. invalid key)
			models = []string{p.ID + "-default", p.ID + "-advanced"}
		}
		
		resp.Data = append(resp.Data, ProviderGroup{
			ProviderID:   p.ID,
			ProviderName: p.Name,
			Models:       models,
		})
	}
	discovery.Mu.RUnlock()

	json.NewEncoder(w).Encode(resp)
}
