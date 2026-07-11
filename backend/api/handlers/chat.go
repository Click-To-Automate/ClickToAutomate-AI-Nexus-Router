package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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

func HandleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Stub for now, can return available models from registry
	fmt.Fprintf(w, `{"object": "list", "data": []}`)
}
