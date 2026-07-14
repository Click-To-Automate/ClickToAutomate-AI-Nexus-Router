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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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

type OpenAIModel struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

type ModelsResponse struct {
	Object string        `json:"object"`
	Data   []OpenAIModel `json:"data"`
}

func HandleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	resp := ModelsResponse{
		Object: "list",
		Data:   []OpenAIModel{},
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
			// Fallback if discovery is pending or failed
			models = []string{p.ID + "-default", p.ID + "-advanced"}
		}
		
		for _, m := range models {
			resp.Data = append(resp.Data, OpenAIModel{
				Id:      m,
				Object:  "model",
				OwnedBy: p.ID, // We inject provider_id here so the frontend can group them
			})
		}
	}
	discovery.Mu.RUnlock()

	json.NewEncoder(w).Encode(resp)
}
