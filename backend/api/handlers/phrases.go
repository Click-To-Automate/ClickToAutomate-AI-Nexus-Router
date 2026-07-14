package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"ainexusrouter-core/apiport"
	"ainexusrouter-core/db"
)

type GeneratePhrasesRequest struct {
	SessionID string `json:"session_id"`
	Text      string `json:"text"`
}

func HandleGetPhrases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id required", http.StatusBadRequest)
		return
	}

	phrases := db.GetActiveLoadingPhrases(sessionID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"phrases": phrases,
	})
}

func HandleGeneratePhrases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req GeneratePhrasesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Respond immediately to avoid blocking the frontend
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"processing"}`))

	// Run generation in background
	go func() {
		// Only generate if text is somewhat meaningful
		if len(req.Text) < 5 {
			return
		}

		prompt := fmt.Sprintf(`You are generating funny, creative, and relatable 'loading' status messages for an AI interface.
Based on the following recent chat history, generate 5 short, witty phrases that show the AI is 'thinking' or 'processing' the request.
Make them sound like fun technical processes or actions related to the user's topic.
They must be formatted strictly as a JSON array of 5 strings. 
Example if the topic is Python: ["Wrangling python snakes...", "Importing BeautifulSoup...", "Indenting with spaces, not tabs...", "Calling the backend...", "Connecting to the mainframe..."]

Recent Chat History:
"""
%s
"""
Output strictly a JSON array of strings and nothing else.`, req.Text)

		chatReq := ChatCompletionRequest{
			Model: "cta-ai-nexus",
			Messages: []map[string]interface{}{
				{"role": "user", "content": prompt},
			},
			Stream: false,
		}

		reqBody, _ := json.Marshal(chatReq)
		
		// Use internal router
		url := fmt.Sprintf("http://127.0.0.1:%s/v1/chat/completions", apiport.Port)
		httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Printf("Generate phrases error: %v", err)
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			log.Printf("Generate phrases error: %v", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		
		var aiResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(body, &aiResp); err != nil || len(aiResp.Choices) == 0 {
			return
		}

		content := aiResp.Choices[0].Message.Content
		
		// Clean up markdown JSON wrappers if present
		content = strings.TrimSpace(content)
		if strings.HasPrefix(content, "```json") {
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)
		} else if strings.HasPrefix(content, "```") {
			content = strings.TrimPrefix(content, "```")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)
		}

		var phrases []string
		if err := json.Unmarshal([]byte(content), &phrases); err != nil {
			log.Printf("Generate phrases parse error: %v\nContent: %s", err, content)
			return
		}

		// Save to DB
		if len(phrases) > 0 {
			db.InsertLoadingPhrases(req.SessionID, phrases)
		}
	}()
}
