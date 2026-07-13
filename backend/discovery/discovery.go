package discovery

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"ainexusrouter-core/config"
	"ainexusrouter-core/db"
)

// isChatModel checks if a model name belongs to an audio, embedding, or moderation model
func isChatModel(id string) bool {
	lower := strings.ToLower(id)
	nonChatKeywords := []string{
		"whisper", "embed", "moderation", "ocr", "tts", "realtime", "prompt-guard",
	}
	for _, kw := range nonChatKeywords {
		if strings.Contains(lower, kw) {
			return false
		}
	}
	return true
}

// ProviderModels stores the discovered live models mapped by providerID
var ProviderModels = make(map[string][]string)
var Mu sync.RWMutex

// RunDiscovery boots up and polls all configured providers to build the model map
func RunDiscovery() {
	log.Println("[Discovery] Initiating Dynamic Model Discovery...")
	
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, p := range config.GlobalConfig.Providers {
		if p.RequiresCustomURL {
			continue // Skip providers requiring custom client-injected URLs
		}
		
		key := db.GetKey(p.ID)
		if key == "" && p.EnvKey != "" {
			key = os.Getenv(p.EnvKey)
		}

		if key == "" {
			continue // Skip providers without keys
		}

		var targetURL string
		// Gemini uses a different path and expects key in query string for models
		if p.ID == "gemini" {
			// Using standard OpenAI compat layer might fail for models endpoint on some APIs,
			// but we will try the /models path from their base_url
			targetURL = p.BaseURL + "/models?key=" + key
		} else {
			targetURL = p.BaseURL + "/models"
		}

		req, err := http.NewRequest(http.MethodGet, targetURL, nil)
		if err != nil {
			log.Printf("[Discovery] Failed to create request for %s: %v\n", p.ID, err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+key)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Discovery] Network error reaching %s: %v\n", p.ID, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var payload map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			log.Printf("[Discovery] Failed to parse models JSON for %s\n", p.ID)
			continue
		}

		var models []string
		
		// Parse OpenAI format: {"data": [{"id": "model-name"}]}
		if data, ok := payload["data"].([]interface{}); ok {
			for _, m := range data {
				if modelObj, ok := m.(map[string]interface{}); ok {
					if id, ok := modelObj["id"].(string); ok {
						if isChatModel(id) {
							models = append(models, id)
						}
					}
				}
			}
		} else if data, ok := payload["models"].([]interface{}); ok {
			// Parse Gemini/Google format: {"models": [{"name": "models/model-name"}]}
			for _, m := range data {
				if modelObj, ok := m.(map[string]interface{}); ok {
					if name, ok := modelObj["name"].(string); ok {
						// Google prepends "models/" to the name
						name = strings.TrimPrefix(name, "models/")
						if isChatModel(name) {
							models = append(models, name)
						}
					}
				}
			}
		}

		if len(models) > 0 {
			Mu.Lock()
			ProviderModels[p.ID] = models
			Mu.Unlock()
			log.Printf("[Discovery] ✅ %s -> Found %d live models", strings.ToUpper(p.ID), len(models))
		}
	}
	
	log.Println("[Discovery] Finished model discovery.")
}

// GetBestModel finds the first discovered model for a provider that contains the keyword.
// If the preferred keyword isn't found, it returns the first available model, or a default string.
func GetBestModel(providerID string, preferredKeyword string, defaultFallback string) string {
	Mu.RLock()
	models, ok := ProviderModels[providerID]
	Mu.RUnlock()
	
	if !ok || len(models) == 0 {
		return defaultFallback // Provider didn't answer discovery, just fallback to hardcoded
	}

	for _, m := range models {
		if strings.Contains(strings.ToLower(m), strings.ToLower(preferredKeyword)) {
			return m
		}
	}

	// Keyword not found, just return the first model we have
	return models[0]
}
