package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"ainexusrouter-core/db"
)

// CompressToToon takes a massive JSON string and uses a fast local API call to Groq 
// to convert it losslessly into Token-Oriented Object Notation (TOON) format.
func CompressToToon(jsonStr string) (string, error) {
	// 1. Fetch the Groq API key from the local database
	apiKey := db.GetKey("groq")
	if apiKey == "" {
		return "", fmt.Errorf("groq API key not found, cannot perform TOON extraction")
	}

	// 2. Build the LLM payload for LLaMA 3 8B
	// We use the 8B model because it's incredibly fast (800+ tokens per second) 
	// and highly capable of formatting data correctly.
	payload := map[string]interface{}{
		"model": "llama-3.1-8b-instant", // Using Groq's fast instant model
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are an expert data compressor. You must convert the user's provided JSON data into TOON (Token-Oriented Object Notation). Use tabular layouts like `[N]{key1,key2}: val1,val2`. Drop unnecessary whitespace and formatting. Do NOT use markdown code blocks like ```toon. Output ONLY the raw TOON text and absolutely no conversational filler or explanations.",
			},
			{
				"role": "user",
				"content": jsonStr,
			},
		},
		"temperature": 0.0, // We want deterministic formatting, zero creativity
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 3. Make synchronous HTTP request to Groq API
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Groq API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// 4. Parse response
	var responseData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(responseData.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from groq")
	}

	result := strings.TrimSpace(responseData.Choices[0].Message.Content)
	
	// Clean up any rogue markdown blocks if the LLM hallucinated them despite system prompt
	result = strings.TrimPrefix(result, "```toon")
	result = strings.TrimPrefix(result, "```TOON")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	return result, nil
}
